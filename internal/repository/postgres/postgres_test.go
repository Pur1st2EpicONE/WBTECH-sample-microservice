package postgres_test

// func TestPostgresStorer_PingClose_Coverage(t *testing.T) {
// 	db, _, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
// 	if err != nil {
// 		t.Fatalf("failed to open sqlmock: %v", err)
// 	}
// 	defer db.Close()
// 	ps := postgres.NewPostgresStorer(sqlx.NewDb(db, "postgres"))
// 	if err := ps.Ping(); err != nil {
// 		t.Errorf("expected Ping success, got %v", err)
// 	}
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
// 	ps := postgres.NewPostgresStorer(db, logger)

// 	order := kafka.CreateOrder(logger)
// 	order.OrderUID = "aboba2"
// 	order.Payment.Transaction = order.OrderUID

// 	if err := ps.SaveOrder(&order); err != nil {
// 		t.Fatalf("SaveOrder failed: %v", err)
// 	}
// }
