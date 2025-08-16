package main

import (
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/app"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/logger"
)

func main() {

	defer logger.CloseFile(logger.OpenFile())
	wbService := app.Start()

	defer wbService.Stop()

	go wbService.RunServer()
	go wbService.RunConsumer()

	wbService.Wait()

}
