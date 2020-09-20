package main

import (
	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/consul"
	"github.com/google/uuid"
	"github.com/hashicorp/consul/api"
	"os"
	"strconv"
)

func Register(consulHost, consulPort, svcHost, svcPort string, logger kitlog.Logger) (register sd.Registrar) {
	var client consul.Client
	{
		consulCfg := api.DefaultConfig()
		consulCfg.Address = consulHost + ":" + consulPort
		consulClient, err := api.NewClient(consulCfg)
		if err != nil {
			logger.Log("create consul client error:", err)
			os.Exit(1)
		}
		client = consul.NewClient(consulClient)
	}

	check := api.AgentServiceCheck{
		HTTP: "http://" + svcHost + ":" + svcPort + "/health",
		Interval: "10s",
		Timeout: "1s",
		Notes: "Consul check service health status.",
	}

	port, _ := strconv.Atoi(svcPort)

	reg := api.AgentServiceRegistration{
		ID: "arithmetic" + uuid.New().String(),
		Name: "arithmetic",
		Address: svcHost,
		Port: port,
		Tags: []string{"arithmetic", "zhengsz"},
		Check: &check,
	}

	register = consul.NewRegistrar(client, &reg, logger)
	return
}
