package repository_test

// func TestConnectDB_Success_Integration(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping integration test in short mode")
// 	}

// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	logger := mock_logger.NewMockLogger(ctrl)
// 	logger.EXPECT().LogInfo("postgres — stopped")

// 	config := configs.Database{
// 		Driver:   "postgres",
// 		Host:     "localhost",
// 		Port:     "5434",
// 		Username: "Neo",
// 		Password: "0451",
// 		DBName:   "wb-service-db-test",
// 		SSLMode:  "disable",
// 	}

// 	db, err := repository.ConnectDB(config)
// 	if err != nil {
// 		t.Fatalf("ConnectPostgres failed: %v", err)
// 	}
// 	defer db.Close()

// 	ps := postgres.NewStorage(db, logger)

// 	if err := ps.Ping(); err != nil {
// 		t.Fatalf("Ping failed: %v", err)
// 	}

// 	ps.Close()
// }

// func TestPostgresStorer_SaveOrder_Integration(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping integration test in short mode")
// 	}
// 	db, err := sqlx.Connect("postgres", "postgres://Neo:0451@localhost:5434/wb-service-db-test?sslmode=disable")
// 	if err != nil {
// 		t.Fatalf("failed to connect to db: %v", err)
// 	}
// 	defer db.Close()

// 	logger := mock_logger.NewMockLogger(gomock.NewController(t))
// 	ps := repository.NewStorage(db, logger)

// 	order := producer.CreateOrder(logger)
// 	order.OrderUID = "1"
// 	order.Payment.Transaction = order.OrderUID

// 	if err := ps.SaveOrder(&order); err != nil {
// 		t.Fatalf("SaveOrder failed: %v", err)
// 	}
// }

// func TestPostgresStorer_GetOrder_Integration(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping integration test in short mode")
// 	}
// 	db, err := sqlx.Connect("postgres", "postgres://Neo:0451@localhost:5434/wb-service-db-test?sslmode=disable")
// 	if err != nil {
// 		t.Fatalf("failed to connect to db: %v", err)
// 	}
// 	defer db.Close()

// 	logger := mock_logger.NewMockLogger(gomock.NewController(t))
// 	ps := repository.NewStorage(db, logger)

// 	order, err := ps.GetOrder("1")
// 	if err != nil {
// 		t.Fatalf("GetOrder failed: %v", err)
// 	}

// 	if order.OrderUID != "1" {
// 		t.Fatalf("expected orderUID 1, got %s", order.OrderUID)
// 	}
// }

// func TestPostgresStorer_SaveAndGetOrder_Integration(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping integration test in short mode")
// 	}
// 	db, err := sqlx.Connect("postgres", "postgres://Neo:0451@localhost:5434/wb-service-db-test?sslmode=disable")
// 	if err != nil {
// 		t.Fatalf("failed to connect to db: %v", err)
// 	}
// 	defer db.Close()
// 	logger := mock_logger.NewMockLogger(gomock.NewController(t))
// 	ps := repository.NewStorage(db, logger)

// 	order := kafka.CreateOrder(logger)

// 	if err := ps.SaveOrder(&order); err != nil {
// 		t.Fatalf("SaveOrder failed: %v", err)
// 	}

// 	gotOrder, err := ps.GetOrder(order.OrderUID)
// 	if err != nil {
// 		t.Fatalf("GetOrder failed: %v", err)
// 	}

// 	if gotOrder.OrderUID != order.OrderUID {
// 		t.Fatalf("expected orderUID %s, got %s", order.OrderUID, gotOrder.OrderUID)
// 	}

// 	if len(gotOrder.Items) != len(order.Items) {
// 		t.Fatalf("expected %d items, got %d", len(order.Items), len(gotOrder.Items))
// 	}
// }
