package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"YAMLGate/config"
	"YAMLGate/internal/gateway"
)

func main() {
	// Получаем путь к конфигурационному файлу из переменной окружения
	configPath := os.Getenv("YAML_GATE_CONFIG")
	if configPath == "" {
		configPath = "config.yaml" // Путь по умолчанию
		fmt.Println("Используется конфигурация по умолчанию:", configPath)
	}

	// Загружаем конфигурацию
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Ошибка при загрузке конфигурации: %v", err)
	}

	// Настраиваем маршруты
	router := gateway.SetupRoutes(cfg)

	// Запускаем сервер
	server := &http.Server{
		Addr:         cfg.Gateway.Address,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	fmt.Printf("Запуск API шлюза на %s\n", cfg.Gateway.Address)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
