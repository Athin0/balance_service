package handlers

import (
	"balance_service/pkg/mErrors"
	"balance_service/pkg/repository"
	"balance_service/pkg/struct4parse"
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Handler struct {
	data   *repository.Repository
	logger *zap.SugaredLogger
}

func NewHandler(data *repository.Repository, logger *zap.SugaredLogger) *Handler {
	return &Handler{
		data:   data,
		logger: logger,
	}
}

func (s *Handler) AddIncome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("{\"errorText\": \"%s\"}", err)))
		log.Println(err)
		return
	}
	incomeParams := &struct4parse.BalanceWithDesc{}
	err = json.Unmarshal(body, incomeParams)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("{\"errorText\": \"%s\"}", err)))
		log.Println(err)
		return
	}
	incomeParams.Time = time.Now()

	err = s.data.AddIncome(r.Context(), *incomeParams)

	if err != nil {
		if errors.Is(err, mErrors.DatabaseError) {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
		w.Write([]byte(fmt.Sprintf("{\"errorText\": \"%s\"}", err)))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{\"status\": \"success\"}"))

}

func (s *Handler) AddReserve(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("{\"errorText\": \"%s\"}", err)))
		return
	}
	incomeParams := &struct4parse.Transaction{}
	err = json.Unmarshal(body, incomeParams)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("{\"errorText\": \"%s\"}", err)))
		return
	}
	incomeParams.Time = time.Now()

	err = s.data.AddReserve(r.Context(), *incomeParams)

	if err != nil {
		if errors.Is(err, mErrors.DatabaseError) {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
		w.Write([]byte(fmt.Sprintf("{\"errorText\": \"%s\"}", err)))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{\"status\": \"success\"}"))

}

func (s *Handler) AddExpense(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("{\"errorText\": \"%s\"}", err)))
		return
	}
	incomeParams := &struct4parse.Transaction{}
	err = json.Unmarshal(body, incomeParams)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("{\"errorText\": \"%s\"}", err)))
		return
	}
	incomeParams.Time = time.Now()

	err = s.data.AddExpense(r.Context(), *incomeParams)

	if err != nil {
		if errors.Is(err, mErrors.DatabaseError) {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
		w.Write([]byte(fmt.Sprintf("{\"errorText\": \"%s\"}", err)))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{\"status\": \"success\"}"))
}

func (s *Handler) GetBalance(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("{\"errorText\": \"%s\"}", err)))
		log.Println(err, "1")
		return
	}
	incomeParams := &struct4parse.Balance{}
	err = json.Unmarshal(body, incomeParams)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("{\"errorText\": \"%s\"}", err)))
		log.Println(err, "2")
		return
	}

	err = s.data.GetBalance(r.Context(), incomeParams)

	if err != nil {
		if errors.Is(err, mErrors.DatabaseError) {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
		w.Write([]byte(fmt.Sprintf("{\"errorText\": \"%s\"}", err)))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{\"status\": \"success\"}"))
	user := strconv.FormatInt(incomeParams.UserId, 10)
	money := strconv.FormatFloat(incomeParams.Value, 'f', -1, 64)
	w.Write([]byte("{\"Balance of user " + user + " : " + money + "}\""))

}

func (s *Handler) GetReserved(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	incomeParams := &[]struct4parse.Transaction{}

	err := s.data.GetAllReserved(r.Context(), incomeParams)

	if err != nil {
		if errors.Is(err, mErrors.DatabaseError) {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
		w.Write([]byte(fmt.Sprintf("{\"errorText\": \"%s\"}", err)))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{\"status\": \"success\"}"))
	fmt.Println(*incomeParams)
	ans := ""
	for _, elem := range *incomeParams {
		ans += "{" + strconv.FormatInt(elem.Id, 10) + ", " + strconv.FormatInt(elem.UserId, 10) + ", " + strconv.FormatInt(elem.ServiceId, 10) + ", " + strconv.FormatInt(elem.OrderId, 10) + ", " + strconv.FormatFloat(elem.Value, 'f', -1, 64) + "},\n"

	}
	w.Write([]byte("{\"Reserves\": [" + ans + "]\"}"))
}

func (s *Handler) GetBalances(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	incomeParams := make([]struct4parse.Balance, 0)

	err := s.data.GetAllBalances(r.Context(), &incomeParams)

	if err != nil {
		if errors.Is(err, mErrors.DatabaseError) {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
		w.Write([]byte(fmt.Sprintf("{\"errorText\": \"%s\"}", err)))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{\"status\": \"success\"}"))
	ans := ""
	for _, elem := range incomeParams {
		ans += "\"" + strconv.FormatInt(elem.UserId, 10) + "\":" + strconv.FormatFloat(elem.Value, 'f', -1, 64) + ", \n"
	}
	w.Write([]byte("{\"Balances\": [" + ans + "]\"}"))
}