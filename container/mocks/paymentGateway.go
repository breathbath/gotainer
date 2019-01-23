package mocks

import (
	"fmt"
)

type PaymentGateway interface {
	CreatePayment(amount float64, clientId string) error
}

//real production gateway
type RealPaymentGateway struct {
	secretKey string
}

func NewRealPaymentGateway(secretKey string) RealPaymentGateway {
	return RealPaymentGateway{secretKey: secretKey}
}

//implements PaymentGateway interface
func (rpg RealPaymentGateway) CreatePayment(amount float64, clientId string) error {
	//here we call the real payment gateway with an http client
	return nil
}

//will handle the registration of users and charge them once a new user is registered
type Registrator struct {
	pg PaymentGateway
}

func NewRegistrator(pg PaymentGateway) Registrator {
	return Registrator{pg}
}

func (r Registrator) RegisterUser(userId string) error {
	//here it registers a user and further charges them 10 USD
	return r.pg.CreatePayment(10, userId)
}

//now I want to test registrator service
type FailingPaymentGateway struct{}

func NewFailingPaymentGateway() FailingPaymentGateway {
	return FailingPaymentGateway{}
}

func (fpg FailingPaymentGateway) CreatePayment(amount float64, clientId string) error {
	return fmt.Errorf("Cannot connect to external api")
}
