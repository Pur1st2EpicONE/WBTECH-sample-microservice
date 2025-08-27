package configs_test

import (
	"os"
	"testing"

	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/configs"
)

func TestLoad_MissingEnvFile(t *testing.T) {
	os.Rename(".env", ".env.aboba")
	defer os.Rename(".env.aboba", ".env")

	_, err := configs.Load()
	if err == nil {
		t.Fatal("expected error due to missing .env, got nil")
	}
}

func TestLoad_MissingConfigFile(t *testing.T) {
	os.Rename("config.yaml", "config.yaml.aboba")
	defer os.Rename("config.yaml.aboba", "config.yaml")

	_, err := configs.Load()
	if err == nil {
		t.Fatal("expected error due to missing config file, got nil")
	}
}

func TestLoad_ValidConfig(t *testing.T) {
	cfg, err := configs.Load()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if cfg.Workers <= 0 {
		t.Errorf("expected workers > 0, got %d", cfg.Workers)
	}
	if cfg.Database.Driver == "" {
		t.Error("expected database driver to be set")
	}
}

func TestProdConfig(t *testing.T) {
	cfg, err := configs.ProdConfig()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	switch {
	case len(cfg.Brokers) == 0:
		t.Error("expected at least one broker to be set")
	case cfg.Topic == "":
		t.Error("expected topic to be set")
	case cfg.ClientID == "":
		t.Error("expected client_id to be set")
	case cfg.MsgsToSend <= 0:
		t.Errorf("expected messages_to_send > 0, got %d", cfg.MsgsToSend)
	}

	k := cfg.Kafka
	if k == nil {
		t.Fatal("expected kafka config to be non-nil")
	}
	switch {
	case k.Acks == "":
		t.Error("expected kafka acks to be set")
	case k.Retries < 0:
		t.Errorf("expected retries >= 0, got %d", k.Retries)
	case k.LingerMs < 0:
		t.Errorf("expected linger_ms >= 0, got %d", k.LingerMs)
	case k.BatchSize <= 0:
		t.Errorf("expected batch_size > 0, got %d", k.BatchSize)
	}
}
