package main

import (
	"balance_service/internal/server"
	"balance_service/pkg/db"
	"balance_service/pkg/repository"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
	"os"
)

func main() {
	zapLogger, errZap := zap.NewProduction()
	if errZap != nil {
		log.Println("Error in creation zapLogger")
	}
	defer func(zapLogger *zap.Logger) {
		err := zapLogger.Sync()
		if err != nil {
			log.Println(err)
		}
	}(zapLogger)
	sugaredLogger := zapLogger.Sugar()

	if err := initConfig(); err != nil {
		sugaredLogger.Fatalf("ошибка инициализации configs: %s", err.Error())
	}

	var err error
	db, err := initDB()
	if err != nil {
		sugaredLogger.Fatalf("ошибка инициализации БД: %s \n", err.Error())
	}

	repo := repository.NewRepository(db)
	hand := server.NewServer(repo, sugaredLogger)
	//r := mux.NewRouter()
	r := fiber.New()

	r.Post("/addIncome", hand.AddIncome)
	r.Post("/addReserve", hand.AddReserve)
	r.Post("/addExpense", hand.AddExpense)
	r.Post("/disReserve", hand.DisReserve)
	r.Get("/getBalance", hand.GetBalance)
	r.Get("/getReserved", hand.GetReserved)
	r.Get("/getBalances", hand.GetBalances)
	r.Get("/getHistory", hand.GetHistory)
	r.Get("/getReport", hand.GetReport)

	r.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))

	r.Use(recover.New())

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("starting service at :" + port)
	err = r.Listen(":" + port)
	if err != nil {
		log.Println("err in listen and serve", err)
		return
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
func initDB() (*db.PostgresDB, error) {
	return db.NewPostgresDB(db.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		User:     viper.GetString("db.user"),
		Password: viper.GetString("db.password"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})
}
