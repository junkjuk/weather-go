package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"weather/internal/domain"
)

type SubscriptionHandler struct {
	service domain.SubscriptionService
}

func NewSubscriptionHandler(service domain.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{service: service}
}

func (s SubscriptionHandler) Subscribe(c *gin.Context) {
	fmt.Println("SubscriptionHandler Subscribe")
	var req domain.SubscribeRequest
	if c.Bind(&req) != nil {
		c.JSON(http.StatusBadRequest, "Invalid input")
		return
	}

	err := s.service.Subscribe(req)
	if err != nil {
		c.JSON(http.StatusNotFound, "Email already subscribed")
		return
	}

	c.JSON(http.StatusOK, "Subscription successful. Confirmation email sent.")
}

func (s SubscriptionHandler) Confirm(c *gin.Context) {
	fmt.Println("SubscriptionHandler Confirm")
	tokenS := c.Param("token")
	token, correct := uuid.Parse(tokenS)

	if correct != nil {
		c.JSON(http.StatusBadRequest, "Invalid token")
		return
	}

	err := s.service.Confirm(token)
	if err != nil {
		c.JSON(http.StatusNotFound, "Token not found")
		return
	}

	c.JSON(http.StatusOK, "Subscription confirmed successfully")
}

func (s SubscriptionHandler) Unsubscribe(c *gin.Context) {
	fmt.Println("SubscriptionHandler Unsubscribe")
	tokenS := c.Param("token")
	token, correct := uuid.Parse(tokenS)

	if correct != nil {
		c.JSON(http.StatusBadRequest, "Invalid token")
		return
	}

	err := s.service.Unsubscribe(token)
	if err != nil {
		c.JSON(http.StatusNotFound, "Token not found")
		return
	}

	c.JSON(http.StatusOK, "Unsubscribed successfully")
}
