package main

import (
	"balance_service/db"
	"balance_service/pkg/handlers"
	"balance_service/pkg/repository"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

var (
	stringLength = viper.GetInt64("uniquestr.len")
	chars        = []rune(viper.GetString("uniquestr.chars"))
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

	stringLength = viper.GetInt64("uniquestr.len")
	chars = []rune(viper.GetString("uniquestr.chars"))

	//var db repository.UrlRepo
	var err error
	db, err := db.NewPostgresDB(db.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		User:     viper.GetString("db.user"),
		Password: viper.GetString("db.password"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})
	if err != nil {
		logger.Fatalf("ошибка инициализации БД: %s \n", err.Error())
	}

	repo := repository.NewRepository(db)
	hand := handlers.NewHandler(repo, logger)
	r := mux.NewRouter()

	r.HandleFunc("/addIncome", hand.AddIncome).Methods("POST")
	r.HandleFunc("/addReserve", hand.AddReserve).Methods("POST")
	r.HandleFunc("/addExpense", hand.AddExpense).Methods("POST")
	r.HandleFunc("/getBalance", hand.GetBalance).Methods("GET")
	r.HandleFunc("/getReserved", hand.GetReserved).Methods("GET")
	r.HandleFunc("/getBalances", hand.GetBalances).Methods("GET")

	fmt.Println("starting service at :" + viper.GetString("port"))
	err = http.ListenAndServe(":"+viper.GetString("port"), r)
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
