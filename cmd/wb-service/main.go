package main

import (
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/app"
)

func main() {

	wbService := app.Start()
	defer wbService.Stop()

	go wbService.RunCacheCleaner()
	go wbService.RunServer()
	go wbService.RunConsumer()

	wbService.Wait()

}
