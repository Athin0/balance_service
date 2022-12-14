package handlers

import (
	"balance_service/pkg/Report"
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
	ans, err := json.Marshal(*incomeParams)
	if err != nil {
		log.Println("err in marshal: ", err)
	}
	w.Write(ans)
}

func (s *Handler) GetReserved(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	incomeParams := &[]struct4parse.Reserve{}

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

	ans, err := json.Marshal(*incomeParams)
	if err != nil {
		log.Println("err in marshal: ", err)
	}
	w.Write(ans)
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
	ans, err := json.Marshal(incomeParams)
	if err != nil {
		log.Println("err in marshal: ", err)
	}
	w.Write(ans)
}

func (s *Handler) GetHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("{\"errorText\": \"%s\"}", err)))
		log.Println(err, "1")
		return
	}
	by := &struct4parse.OrderParams{}
	err = json.Unmarshal(body, by)
	if err != nil && string(body) != "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("{\"errorText\": \"%s\"}", err)))
		log.Println(err, "2")
		return
	}
	incomeParams := make([]struct4parse.Transaction, 0)
	err = s.data.GetAllTransactions(r.Context(), &incomeParams, *by)
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

	since := 0
	if by.Since >= len(incomeParams) {
		return
	} else {
		since = by.Since
	}
	until := len(incomeParams)
	if by.Num != 0 && by.Since+by.Num < len(incomeParams) {
		until = by.Since + by.Num
	}
	incomeParams = incomeParams[since:until]
	ans, err := json.Marshal(incomeParams)
	if err != nil {
		log.Println("err in marshal: ", err)
	}
	w.Write(ans)
}

func (s *Handler) DisReserve(w http.ResponseWriter, r *http.Request) {
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

	err = s.data.DisReserve(r.Context(), *incomeParams)

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

func (s *Handler) GetReport(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("{\"errorText\": \"%s\"}", err)))
		log.Println(err)
		return
	}
	timeDur := &struct4parse.Time4Report{}
	err = json.Unmarshal(body, timeDur)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("{\"errorText\": \"%s\"}", err)))
		log.Println(err)
		return
	}

	incomeParams := make([]struct4parse.Report, 0)

	err = s.data.GetReports(r.Context(), &incomeParams, *timeDur)
	text, err := Report.MakeReport(&incomeParams)
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
	w.Write([]byte("{\"file name\": \"" + text + "\"}"))
	ans, err := json.Marshal(incomeParams)
	if err != nil {
		log.Println("err in marshal: ", err)
	}
	w.Write(ans)
}
