package main

import (
	"github.com/nmluci/cache-optimization/cmd/webservice"
	"github.com/nmluci/cache-optimization/internal/component"
	"github.com/nmluci/cache-optimization/internal/config"
)

func main() {
	config.Init()
	conf := config.Get()

	logger := component.NewLogger(component.NewLoggerParams{
		ServiceName: conf.ServiceName,
		PrettyPrint: true,
	})

	webservice.Start(conf, logger)
}
