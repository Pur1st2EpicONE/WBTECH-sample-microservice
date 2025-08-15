package main

import (
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/app"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/logger"
)

func main() {

	defer logger.CloseFile(logger.OpenFile())
	wbService := app.New()

	ctx, cancel := wbService.NewContext()
	defer cancel()

	go wbService.RunServer(ctx)
	go wbService.RunConsumer(ctx)

	wbService.Wait(ctx)

}
