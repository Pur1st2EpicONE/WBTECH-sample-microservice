package config

import "github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/repository/postgres"

var Pg = postgres.PgConfig{
	Host:     "localhost",
	Port:     "5433",
	Username: "Neo",
	Password: "0451",
	DBName:   "wb-service-db",
	SSLMode:  "disable",
}

var HTTPPort = "8080"
var KafkaBrokers = []string{"localhost:9092"}
var KafkaTopic = "orders"
