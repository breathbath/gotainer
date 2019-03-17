package main

import (
	"example/cont"
	"example/email"
	"github.com/breathbath/gotainer/container"
)

func main() {
	builder := container.RuntimeContainerBuilder{}
	runtimeContainer, err := builder.BuildContainerFromConfig(cont.GetConfig())
	if err != nil {
		panic(err)
	}

	smtpClient := runtimeContainer.Get("smtpClient", true).(email.SmtpClient)
	err = smtpClient.SendEmail("no@mail.me", "Client", "Hey", "Hello world")
	if err != nil {
		panic(err)
	}
}
