package handlers4HttpRouter

import (
	"balance_service/pkg/mErrors"
	"balance_service/pkg/repository"
	"balance_service/pkg/utils"
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Server struct {
	data   *repository.Repository
	logger *zap.SugaredLogger
}

func NewServer(data *repository.Repository, logger *zap.SugaredLogger) *Server {
	return &Server{
		data:   data,
		logger: logger,
	}
}

func (s *Server) AddIncome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("{\"errorText\": \"%s\"}", err)))
		log.Println(err)
		return
	}
	incomeParams := &utils.BalanceWithDesc{}
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

func (s *Server) AddReserve(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("{\"errorText\": \"%s\"}", err)))
		return
	}
	incomeParams := &utils.Transaction{}
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

func (s *Server) AddExpense(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("{\"errorText\": \"%s\"}", err)))
		return
	}
	incomeParams := &utils.Transaction{}
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

func (s *Server) GetBalance(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("{\"errorText\": \"%s\"}", err)))
		log.Println(err, "1")
		return
	}
	incomeParams := &utils.Balance{}
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

func (s *Server) GetReserved(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	incomeParams := &[]utils.Reserve{}

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

func (s *Server) GetBalances(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	incomeParams := make([]utils.Balance, 0)

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

func (s *Server) GetHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("{\"errorText\": \"%s\"}", err)))
		log.Println(err, "1")
		return
	}
	by := &utils.OrderParams{}
	err = json.Unmarshal(body, by)
	if err != nil && string(body) != "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("{\"errorText\": \"%s\"}", err)))
		log.Println(err, "2")
		return
	}
	incomeParams := make([]utils.Transaction, 0)
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

func (s *Server) DisReserve(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("{\"errorText\": \"%s\"}", err)))
		return
	}
	incomeParams := &utils.Transaction{}
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

func (s *Server) GetReport(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("{\"errorText\": \"%s\"}", err)))
		log.Println(err)
		return
	}
	timeDur := &utils.Time4Report{}
	err = json.Unmarshal(body, timeDur)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("{\"errorText\": \"%s\"}", err)))
		log.Println(err)
		return
	}

	incomeParams := make([]utils.Report, 0)

	err = s.data.GetReports(r.Context(), &incomeParams, *timeDur)
	text, err := utils.MakeReport(&incomeParams)
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
