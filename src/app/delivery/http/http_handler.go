package http

import (
	"context"
	"log"
	"net/http"
	app "work-relic-bot/src/app/usecase"
	"work-relic-bot/src/models"

	validator "gopkg.in/go-playground/validator.v9"

	"github.com/gin-gonic/gin"
)

type ResponseError struct {
	Message string `json:"message"`
}

type HttpAppHandler struct {
	AUsecase app.IUsecase
}

func NewHttpHandler(r *gin.Engine, au app.IUsecase) {
	handler := &HttpAppHandler{
		AUsecase: au,
	}
	r.GET("/check", handler.HealthCheck)
	r.POST("/webhook", handler.Webhook)
}

func validateRequest(m *models.Webhook) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (a *HttpAppHandler) Webhook(c *gin.Context) {
	var webhook models.Webhook
	err := c.Bind(&webhook)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	if ok, err := validateRequest(&webhook); !ok {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	ctx := c.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	res, err := a.AUsecase.Webhook(ctx, &webhook)

	if err != nil {
		c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		return
	}
	c.Header("Content-Type", "applcation/json")
	c.JSON(http.StatusOK, gin.H{
		"status": res,
	})
	return
}

func (a *HttpAppHandler) HealthCheck(c *gin.Context) {
	ctx := c.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	res, err := a.AUsecase.Check(ctx)

	if err != nil {
		c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		return
	}
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, gin.H{
		"status": res,
	})
	return
}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}
	log.Println(err)
	switch err {
	case models.ErrInternalServerError:
		return http.StatusInternalServerError
	case models.ErrNotFound:
		return http.StatusNotFound
	case models.ErrConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
