package main

import (
	"errors"
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/go-kit/kit/log"
	"github.com/hashicorp/consul/api"
	"github.com/openzipkin/zipkin-go"
	zipkinhttpsvr "github.com/openzipkin/zipkin-go/middleware/http"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"strings"
	"sync"
)

type HystrixRouter struct {
	svcMap       *sync.Map
	logger       log.Logger
	fallbackMsg  string
	consulClient *api.Client
	tracer       *zipkin.Tracer
}

func Routes(client *api.Client, zipkinTracer *zipkin.Tracer, fbMsg string, logger log.Logger) http.Handler {
	return HystrixRouter{
		svcMap:       &sync.Map{},
		logger:       logger,
		fallbackMsg:  fbMsg,
		consulClient: client,
		tracer:       zipkinTracer,
	}
}
func (router HystrixRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 查询原始请求路径： /arithmetic/calculate/10/5
	reqPath := r.URL.Path
	if reqPath == "" {
		return
	}

	// 分割获取serviceName
	pathArray := strings.Split(reqPath, "/")
	serviceName := pathArray[1]

	// 检查是否已经加入监控
	if _, ok := router.svcMap.Load(serviceName); !ok {
		// 将serviceName作为命令对象，设置参数
		hystrix.ConfigureCommand(serviceName, hystrix.CommandConfig{Timeout: 1000})
		router.svcMap.Store(serviceName, serviceName)
	}
	// 执行命令
	err := hystrix.Do(serviceName, func() (err error) {
		result, _, err := router.consulClient.Catalog().Service(serviceName, "", nil)
		if err != nil {
			router.logger.Log("ReverseProxy failed", "query service instance error", err.Error())
			return
		}

		if len(result) == 0 {
			router.logger.Log("ReverseProxy failed", "query service instance error", err.Error())
			return errors.New("no such service instance")
		}

		director := func(req *http.Request) {
			// 重新组织请求路径，去掉服务名称部分
			destPath := strings.Join(pathArray[2:], "/")

			// 随机选择一个服务实例
			tgt := result[rand.Int()%len(result)]
			router.logger.Log("service id", tgt.ServiceID)

			// 设置代理服务地址信息
			req.URL.Scheme = "http"
			req.URL.Host = fmt.Sprintf("%s:%d", tgt.ServiceAddress, tgt.ServicePort)
			req.URL.Path = "/" + destPath
		}

		var proxyError error = nil
		roundTrip, _ := zipkinhttpsvr.NewTransport(router.tracer, zipkinhttpsvr.TransportTrace(true))

		errorHandler := func(ew http.ResponseWriter, er *http.Request, err error) {
			proxyError = err
		}

		proxy := &httputil.ReverseProxy{
			Director:     director,
			Transport:    roundTrip,
			ErrorHandler: errorHandler,
		}
		proxy.ServeHTTP(w, r)
		return proxyError
	}, func(err error) error {
		router.logger.Log("fallback error description", err.Error())
		return errors.New(router.fallbackMsg)
	})
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
	}
}
