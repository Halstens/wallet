// package main

// import "flag"

// func main() {

// 	// addr := flag.String("addr", ":4000", "Сетевой адрес HTTP")
// 	// dsn := flag.String("dsn", "%s:%s@/wallet?parseTime=true", "walletdb", POSTGRES_USER, POSTGRES_PASSWORD)

// }
package main

import (
	"log"

	"github.com/wallet/internal/config"
	"github.com/wallet/internal/repository"
)

func main() {
	// Загрузка конфига
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Подключение к БД
	db, err := repository.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("Successfully connected to PostgreSQL!")
	// Далее инициализация сервера...
}
