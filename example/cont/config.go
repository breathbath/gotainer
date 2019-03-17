package cont

import (
	"github.com/breathbath/gotainer/container"
	"example/email"
	"example/passwords"
)

func GetConfig() container.Tree {
	return container.Tree{
		container.Node{
			Parameters: map[string]interface{}{
				"fromEmail": "admin@mail.me",
				"fromName":   "Admin",
			},
		},
		container.Node{
			ID:           "passwordManager",
			NewFunc:      passwords.NewPasswordManager,
		},
		container.Node{
			ID:           "smtpClient",
			NewFunc:      email.NewSmtpClient,
			ServiceNames: container.Services{
				"passwordManager", //with this we say to inject passwordManager into email.NewSmtpClient method
			},
		},
	}
}
