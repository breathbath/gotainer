package email

import "fmt"

//we abstract the login and password provider
type SmtpLoginPasswordProvider interface {
	GetLoginData() (login, pass string, err error)
}

type SmtpClient struct {
	loginPassProvider SmtpLoginPasswordProvider
	fromEmail         string
	fromName          string
}

//constructor
func NewSmtpClient(loginPassProvider SmtpLoginPasswordProvider, fromEmail, fromName string) SmtpClient {
	return SmtpClient{
		loginPassProvider: loginPassProvider,
		fromEmail:         fromEmail,
		fromName:          fromName,
	}
}

func (smtp SmtpClient) SendEmail(receiverEmail, receiverName, subject, body string) error {
	login, pass, err := smtp.loginPassProvider.GetLoginData()
	if err != nil {
		return err
	}

	fmt.Println(login, pass)
	//we do here the actual logic for sending emails via smtp

	return nil
}
