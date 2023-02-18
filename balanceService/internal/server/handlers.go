package server

import (
	"balance_service/pkg/mErrors"
	"balance_service/pkg/repository"
	"balance_service/pkg/utils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
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

func (s *Server) AddIncome(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/json")
	incomeParams := &utils.BalanceWithDesc{}
	err := c.BodyParser(incomeParams)
	if err != nil {
		c.SendStatus(http.StatusBadRequest)
		c.Write([]byte(fmt.Sprintf("{\"errorText\": \"%s\"}", err)))
		return err
	}
	incomeParams.Time = time.Now()
	err = s.data.AddIncome(c.Context(), *incomeParams)

	if err != nil {
		if errors.Is(err, mErrors.DatabaseError) {
			c.SendStatus(http.StatusInternalServerError)
		} else {
			c.SendStatus(http.StatusBadRequest)
		}
		c.Write([]byte(fmt.Sprintf("{\"errorText\": \"%s\"}", err)))
		return err
	}
	c.SendStatus(http.StatusOK)
	c.Write([]byte("{\"status\": \"success\"}"))
	return nil
}

func (s *Server) AddReserve(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/json")
	incomeParams := &utils.Transaction{}
	err := c.BodyParser(incomeParams)
	if err != nil {
		c.SendStatus(http.StatusBadRequest)
		c.Write([]byte(fmt.Sprintf("{\"errorText\": \"%s\"}", err)))
		return err
	}
	incomeParams.Time = time.Now()

	err = s.data.AddReserve(c.Context(), *incomeParams)

	if err != nil {
		if errors.Is(err, mErrors.DatabaseError) {
			c.SendStatus(http.StatusInternalServerError)
		} else {
			c.SendStatus(http.StatusBadRequest)
		}
		c.Write([]byte(fmt.Sprintf("{\"errorText\": \"%s\"}", err)))
		return err
	}
	c.SendStatus(http.StatusOK)
	c.Write([]byte("{\"status\": \"success\"}"))
	return nil
}

func (s *Server) AddExpense(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/json")

	incomeParams := &utils.Transaction{}
	err := c.BodyParser(incomeParams)
	if err != nil {
		c.SendStatus(http.StatusBadRequest)
		c.Write([]byte(fmt.Sprintf("{\"errorText\": \"%s\"}", err)))
		return err
	}
	incomeParams.Time = time.Now()

	err = s.data.AddExpense(c.Context(), *incomeParams)

	if err != nil {
		if errors.Is(err, mErrors.DatabaseError) {
			c.SendStatus(http.StatusInternalServerError)
		} else {
			c.SendStatus(http.StatusBadRequest)
		}
		c.Write([]byte(fmt.Sprintf("{\"errorText\": \"%s\"}", err)))
		return err
	}
	c.SendStatus(http.StatusOK)
	c.Write([]byte("{\"status\": \"success\"}"))
	return nil
}

func (s *Server) GetBalance(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/json")
	incomeParams := &utils.Balance{}
	err := c.BodyParser(incomeParams)
	if err != nil {
		c.SendStatus(http.StatusBadRequest)
		c.Write([]byte(fmt.Sprintf("{\"errorText\": \"%s\"}", err)))

		return err
	}

	err = s.data.GetBalance(c.Context(), incomeParams)

	if err != nil {
		if errors.Is(err, mErrors.DatabaseError) {
			c.SendStatus(http.StatusInternalServerError)
		} else {
			c.SendStatus(http.StatusBadRequest)
		}
		c.Write([]byte(fmt.Sprintf("{\"errorText\": \"%s\"}", err)))
		return err
	}
	c.SendStatus(http.StatusOK)
	c.Write([]byte("{\"status\": \"success\"}"))
	ans, err := json.Marshal(*incomeParams)
	if err != nil {
		log.Println("err in marshal: ", err)
	}
	c.Write(ans)
	return nil
}

func (s *Server) GetReserved(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/json")

	incomeParams := &[]utils.Reserve{}

	err := s.data.GetAllReserved(c.Context(), incomeParams)

	if err != nil {
		if errors.Is(err, mErrors.DatabaseError) {
			c.SendStatus(http.StatusInternalServerError)
		} else {
			c.SendStatus(http.StatusBadRequest)
		}
		c.Write([]byte(fmt.Sprintf("{\"errorText\": \"%s\"}", err)))
		return err
	}
	c.SendStatus(http.StatusOK)
	c.Write([]byte("{\"status\": \"success\"}"))

	ans, err := json.Marshal(*incomeParams)
	if err != nil {
		log.Println("err in marshal: ", err)
	}
	c.Write(ans)
	return nil
}

func (s *Server) GetBalances(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/json")

	incomeParams := make([]utils.Balance, 0)

	err := s.data.GetAllBalances(c.Context(), &incomeParams)

	if err != nil {
		if errors.Is(err, mErrors.DatabaseError) {
			c.SendStatus(http.StatusInternalServerError)
		} else {
			c.SendStatus(http.StatusBadRequest)
		}
		c.Write([]byte(fmt.Sprintf("{\"errorText\": \"%s\"}", err)))
		return err
	}
	c.SendStatus(http.StatusOK)
	c.Write([]byte("{\"status\": \"success\"}"))
	ans, err := json.Marshal(incomeParams)
	if err != nil {
		log.Println("err in marshal: ", err)
	}
	c.Write(ans)
	return nil
}

func (s *Server) GetHistory(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/json")
	by := &utils.OrderParams{}
	err := c.BodyParser(by)
	if err != nil {
		c.SendStatus(http.StatusBadRequest)
		c.Write([]byte(fmt.Sprintf("{\"errorText\": \"%s\"}", err)))
		return err
	}
	if err != nil && len(c.Body()) != 0 { //""- empty request
		c.SendStatus(http.StatusBadRequest)
		c.Write([]byte(fmt.Sprintf("{\"errorText\": \"%s\"}", err)))
		return err
	}
	incomeParams := make([]utils.Transaction, 0)
	err = s.data.GetAllTransactions(c.Context(), &incomeParams, *by)
	if err != nil {
		if errors.Is(err, mErrors.DatabaseError) {
			c.SendStatus(http.StatusInternalServerError)
		} else {
			c.SendStatus(http.StatusBadRequest)
		}
		c.Write([]byte(fmt.Sprintf("{\"errorText\": \"%s\"}", err)))
		return err
	}
	c.SendStatus(http.StatusOK)
	c.Write([]byte("{\"status\": \"success\"}"))

	since := 0
	if by.Since >= len(incomeParams) {
		return err
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
	c.Write(ans)
	return nil
}

func (s *Server) DisReserve(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/json")

	incomeParams := &utils.Transaction{}
	err := c.BodyParser(incomeParams)
	if err != nil {
		c.SendStatus(http.StatusBadRequest)
		c.Write([]byte(fmt.Sprintf("{\"errorText\": \"%s\"}", err)))
		return err
	}
	incomeParams.Time = time.Now()

	err = s.data.DisReserve(c.Context(), *incomeParams)

	if err != nil {
		if errors.Is(err, mErrors.DatabaseError) {
			c.SendStatus(http.StatusInternalServerError)
		} else {
			c.SendStatus(http.StatusBadRequest)
		}
		c.Write([]byte(fmt.Sprintf("{\"errorText\": \"%s\"}", err)))
		return err
	}
	c.SendStatus(http.StatusOK)
	c.Write([]byte("{\"status\": \"success\"}"))
	return nil
}

func (s *Server) GetReport(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/json")
	timeDur := &utils.Time4Report{}
	err := c.BodyParser(timeDur)
	if err != nil {
		c.SendStatus(http.StatusBadRequest)
		c.Write([]byte(fmt.Sprintf("{\"errorText\": \"%s\"}", err)))

		return err
	}

	incomeParams := make([]utils.Report, 0)

	err = s.data.GetReports(c.Context(), &incomeParams, *timeDur)
	text, err := utils.MakeReport(&incomeParams)
	if err != nil {
		if errors.Is(err, mErrors.DatabaseError) {
			c.SendStatus(http.StatusInternalServerError)
		} else {
			c.SendStatus(http.StatusBadRequest)
		}
		c.Write([]byte(fmt.Sprintf("{\"errorText\": \"%s\"}", err)))
		return err
	}
	c.SendStatus(http.StatusOK)
	c.Write([]byte("{\"status\": \"success\"}"))
	c.Write([]byte("{\"file name\": \"" + text + "\"}"))
	ans, err := json.Marshal(incomeParams)
	if err != nil {
		log.Println("err in marshal: ", err)
	}
	c.Write(ans)
	return nil
}
