package cont

import (
	"github.com/breathbath/gotainer/container"
	"github.com/breathbath/gotainer/container/example/email"
	"github.com/breathbath/gotainer/container/example/passwords"
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
				"fromEmail",
				"fromName",
			},
		},
	}
}
