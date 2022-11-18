package main

import (
	"balance_service/db"
	"balance_service/pkg/handlers"
	"balance_service/pkg/middleware"
	"balance_service/pkg/repository"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
	"net/http"
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
	logger := zapLogger.Sugar()

	if err := initConfig(); err != nil {
		logger.Fatalf("ошибка инициализации configs: %s", err.Error())
	}

	var err error
	db, err := initDB()
	if err != nil {
		logger.Fatalf("ошибка инициализации БД: %s \n", err.Error())
	}

	repo := repository.NewRepository(db)
	hand := handlers.NewHandler(repo, logger)
	r := mux.NewRouter()

	r.HandleFunc("/addIncome", hand.AddIncome).Methods("POST")
	r.HandleFunc("/addReserve", hand.AddReserve).Methods("POST")
	r.HandleFunc("/addExpense", hand.AddExpense).Methods("POST")
	r.HandleFunc("/disReserve", hand.DisReserve).Methods("POST")
	r.HandleFunc("/getBalance", hand.GetBalance).Methods("GET")
	r.HandleFunc("/getReserved", hand.GetReserved).Methods("GET")
	r.HandleFunc("/getBalances", hand.GetBalances).Methods("GET")
	r.HandleFunc("/getHistory", hand.GetHistory).Methods("GET")
	r.HandleFunc("/getReport", hand.GetReport).Methods("GET")

	r0 := middleware.AccessLog(logger, r)
	r0 = middleware.Panic(r0)

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("starting service at :" + port)
	err = http.ListenAndServe(":"+port, r0)
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
