package main

import (
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/app"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/logger"
)

func main() {

	defer logger.CloseFile(logger.OpenFile())

	orderService := app.New()

	go orderService.RunServer()
	go orderService.RunConsumer()

	orderService.WaitForShutdown()

}
