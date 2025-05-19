package domain

import (
	"fmt"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
	"weather/configs"
)

type SubscriptionRepo struct {
	db  *gorm.DB
	cfg *configs.Config
}

func NewSubscriptionRepo(cfg *configs.Config) *SubscriptionRepo {

	dsn := fmt.Sprintf("host=postgres user=%s password=%s dbname=%s port=5432", cfg.DBUser, cfg.DBPass, cfg.DBName)
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database!")
	}

	database.AutoMigrate(&Subscription{})

	return &SubscriptionRepo{
		db: database,
	}
}

func (repo *SubscriptionRepo) Add(subscription *Subscription) {
	repo.db.Create(subscription)
}

func (repo *SubscriptionRepo) GetAll() []Subscription {
	var subscriptions []Subscription
	repo.db.Find(&subscriptions)
	return subscriptions
}

func (repo *SubscriptionRepo) GetAllActive() []Subscription {
	var subscriptions []Subscription
	repo.db.Where("is_active = ?", true).Find(&subscriptions)
	return subscriptions
}

func (repo *SubscriptionRepo) GetById(id uuid.UUID) Subscription {
	var subscription Subscription
	repo.db.First(&subscription, "ID = ?", id)
	return subscription
}

func (repo *SubscriptionRepo) GetByEmail(email string) Subscription {
	var subscription Subscription
	repo.db.Where("is_active = ?", true).First(&subscription, "email = ?", email)
	return subscription
}

func (repo *SubscriptionRepo) SetActive(id uuid.UUID) {
	var subscription Subscription
	repo.db.First(&subscription, "ID = ?", id)
	repo.db.Model(&subscription).Update("IsActive", true)
}

func (repo *SubscriptionRepo) SetUnactive(id uuid.UUID) {
	var subscription Subscription
	repo.db.First(&subscription, "ID = ?", id)
	repo.db.Model(&subscription).Update("IsActive", false)
}

func (repo *SubscriptionRepo) UpdateLastRun(id uuid.UUID) {
	repo.db.Model(&Subscription{}).Where("ID = ?", id).Update("LastRun", time.Now())
}
