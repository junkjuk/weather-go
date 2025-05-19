package domain

import (
	"errors"
	"github.com/google/uuid"
)

type SubscriptionService interface {
	Subscribe(subscribe SubscribeRequest) error
	Confirm(token uuid.UUID) error
	Unsubscribe(token uuid.UUID) error
}

type SubscriptionServiceImpl struct {
	db     *SubscriptionRepo
	sender *EmailSender
}

func NewSubscriptionServiceImpl(db *SubscriptionRepo, sender *EmailSender) SubscriptionServiceImpl {
	return SubscriptionServiceImpl{
		db:     db,
		sender: sender,
	}
}

func (s SubscriptionServiceImpl) Subscribe(subscribe SubscribeRequest) error {

	subscription := s.db.GetByEmail(subscribe.Email)
	if subscription.Email != "" {
		return errors.New("already subscribed")
	}

	token := uuid.New()

	newSubscription := Subscription{
		ID:        token,
		Email:     subscribe.Email,
		City:      subscribe.City,
		Frequency: subscribe.Frequency,
		IsActive:  false,
	}
	s.db.Add(&newSubscription)
	mail := GetConfirmEmail(newSubscription)
	s.sender.SendEmailHtml(subscribe.Email, "New subscription", mail)
	return nil
}

func (s SubscriptionServiceImpl) Confirm(token uuid.UUID) error {
	sub := s.db.GetById(token)
	if sub.ID != token {
		return errors.New("token not found")
	}
	s.db.SetActive(token)

	return nil
}

func (s SubscriptionServiceImpl) Unsubscribe(token uuid.UUID) error {
	sub := s.db.GetById(token)
	if sub.ID != token {
		return errors.New("token not found")
	}
	s.db.SetUnactive(token)
	return nil
}
