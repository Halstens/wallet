package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/wallet/internal/database/postgress"

	_ "github.com/lib/pq"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	wallets  *postgress.WalletRepository
}

func main() {
	addr := flag.String("addr", ":4000", "Сетевой адрес HTTP")
	//dsn := flag.String("dsn", "Pavel:fhaar355228F@/tralaleo?parseTime=true", "tralaleo_db")
	dsn := flag.String("dsn", "user=postgres password=fhaar355228F dbname=walletdb sslmode=disable", "walletdb")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		wallets:  &postgress.WalletRepository{DB: db},
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Запуск сервера на %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

// // package main

// // import "flag"

// // func main() {

// // 	// addr := flag.String("addr", ":4000", "Сетевой адрес HTTP")
// // 	// dsn := flag.String("dsn", "%s:%s@/wallet?parseTime=true", "walletdb", POSTGRES_USER, POSTGRES_PASSWORD)

// // }
// package main

// import (
// 	"log"

// 	"github.com/wallet/internal/config"
// 	"github.com/wallet/internal/repository"
// )

// func main() {
// 	// Загрузка конфига
// 	cfg, err := config.Load()
// 	if err != nil {
// 		log.Fatalf("Failed to load config: %v", err)
// 	}

// 	// Подключение к БД
// 	db, err := repository.NewPostgresDB(cfg)
// 	if err != nil {
// 		log.Fatalf("Failed to connect to database: %v", err)
// 	}
// 	defer db.Close()

// 	log.Println("Successfully connected to PostgreSQL!")
// 	// Далее инициализация сервера...
// }
