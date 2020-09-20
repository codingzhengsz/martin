package main

import (
	"flag"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/hashicorp/consul/api"
	"github.com/openzipkin/zipkin-go"
	zipkinhttpsvr "github.com/openzipkin/zipkin-go/middleware/http"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

// 需要在代码中嵌入Reporter组件

func NewReverseProxy(client *api.Client, zipkinTracer *zipkin.Tracer, logger log.Logger) *httputil.ReverseProxy {
	director := func(req *http.Request) {
		reqPath := req.URL.Path
		if reqPath == "" {
			return
		}
		pathArray := strings.Split(reqPath, "/")
		serviceName := pathArray[1]

		result, _, err := client.Catalog().Service(serviceName, "", nil)
		if err != nil {
			logger.Log("ReverseProxy failed", "query service instance error", err.Error())
			return
		}

		if len(result) == 0 {
			logger.Log("ReverseProxy failed", "no such service instance", serviceName)
			return
		}

		destPath := strings.Join(pathArray[2:], "/")

		tgt := result[rand.Int()%len(result)]
		logger.Log("service id", tgt.ServiceID)

		// 设置代理服务地址信息
		req.URL.Scheme = "http"
		req.URL.Host = fmt.Sprintf("%s:%d", tgt.ServiceAddress, tgt.ServicePort)
		req.URL.Path = "/" + destPath

	}
	roundTrip, _ := zipkinhttpsvr.NewTransport(zipkinTracer, zipkinhttpsvr.TransportTrace(true))
	return &httputil.ReverseProxy{Director: director, Transport: roundTrip}
}

func main() {
	var (
		consulHost = flag.String("consul.host", "localhost", "consul server ip address")
		consulPort = flag.String("consul.port", "8500", "consul server port")
		zipkinUrl  = flag.String("zipkin.url", "http://localhost:9411/api/v2/spans", "Zipkin server url")
	)
	flag.Parse()

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	var zipkinTracer *zipkin.Tracer
	{
		var (
			err           error
			hostPort      = "localhost:9090"
			serviceName   = "gateway-service"
			useNoopTracer = *zipkinUrl == ""
			reporter      = zipkinhttp.NewReporter(*zipkinUrl)
		)
		defer reporter.Close()
		zEP, _ := zipkin.NewEndpoint(serviceName, hostPort)
		zipkinTracer, err = zipkin.NewTracer(
			reporter, zipkin.WithLocalEndpoint(zEP), zipkin.WithNoopTracer(useNoopTracer))
		if err != nil {
			logger.Log("err", err)
			os.Exit(1)
		}
		if !useNoopTracer {
			logger.Log("tracer", "Zipkin", "type", "Native", "URL", *zipkinUrl)
		}
	}

	// 创建consul api客户端
	consulConfig := api.DefaultConfig()
	consulConfig.Address = "http://" + *consulHost + ":" + *consulPort
	consulClient, err := api.NewClient(consulConfig)
	if err != nil {
		logger.Log("err", err)
		os.Exit(1)
	}

	//创建反向代理
	proxy := NewReverseProxy(consulClient, zipkinTracer, logger)

	tags := map[string]string{
		"component": "gateway_server",
	}

	handler := zipkinhttpsvr.NewServerMiddleware(
		zipkinTracer,
		zipkinhttpsvr.SpanName("gateway"),
		zipkinhttpsvr.TagResponseSize(true),
		zipkinhttpsvr.ServerTags(tags),
	)(proxy)

	errc := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()

	//开始监听
	go func() {
		logger.Log("transport", "HTTP", "addr", "9090")
		errc <- http.ListenAndServe(":9090", handler)
	}()

	// 开始运行，等待结束
	logger.Log("exit", <-errc)
}
