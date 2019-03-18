package main

import (
	"github.com/breathbath/gotainer/container"
	"github.com/breathbath/gotainer/container/example/cont"
	"github.com/breathbath/gotainer/container/example/email"
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
