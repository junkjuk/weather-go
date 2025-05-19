package domain

import (
	"github.com/google/uuid"
	"time"
)

type SubscribeRequest struct {
	Email     string `json:"email"`
	City      string `json:"city"`
	Frequency string `json:"frequency" binding:"oneof=hourly daily"`
}

type Subscription struct {
	ID        uuid.UUID `gorm:"primaryKey"`
	Email     string
	City      string
	Frequency string `gorm:"type:varchar(10); check (frequency IN ('hourly', 'daily'))"`
	IsActive  bool
	LastRun   time.Time
}
