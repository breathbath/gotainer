package container

import (
	"errors"
	"fmt"
)

// This example scanning of services from container with a possible error output
// if it happens
func ExampleRuntimeContainer_ScanSecure() {
	cont := NewRuntimeContainer()
	cont.AddConstructor("parameterOk", func(c Container) (interface{}, error) {
		return "123456", nil
	})
	cont.AddConstructor("parameterFailed", func(c Container) (interface{}, error) {
		return "123456", errors.New("Some error")
	})

	var parameterOk string
	err := cont.ScanSecure("parameterOk", true, &parameterOk)
	fmt.Println(err)

	err = cont.ScanSecure("parameterFailed", true, &parameterOk)
	fmt.Println(err)

	err = cont.ScanSecure("unknownParameter", true, &parameterOk)
	fmt.Println(err)

	// Output:
	// <nil>
	// Some error [check 'parameterFailed' service]
	// Unknown dependency 'unknownParameter'
}

// This example shows how to add a custom constructor and use an existing service in it
func ExampleRuntimeContainer_AddConstructor() {
	cont := NewRuntimeContainer()
	cont.AddConstructor("secretKey", func(c Container) (interface{}, error) {
		return "00000", nil
	})
	cont.AddConstructor("login", func(c Container) (interface{}, error) {
		return "admin", nil
	})

	cont.AddConstructor("credentials", func(c Container) (interface{}, error) {
		secretKey := c.Get("secretKey", true).(string)
		login := c.Get("login", true).(string)

		return fmt.Sprintf("%s:%s", login, secretKey), nil
	})

	credentials := cont.Get("credentials", true).(string)
	fmt.Println(credentials)

	// Output:
	// admin:00000
}

type Db struct{}

func NewDb() Db {
	return Db{}
}

type NameProvider struct {
	db Db
}

func (db Db) GetName() string {
	return "Elton John"
}

func NewNameProvider(db Db) NameProvider {
	return NameProvider{db: db}
}

func (dp NameProvider) GetFullName(salutation string) string {
	return salutation + " " + dp.db.GetName()
}

// This example shows how to add a struct constructor with some dependencies
func ExampleRuntimeContainer_AddNewMethod() {
	//having declared struct dependencies:
	//type Db struct{}
	//
	//func NewDb() Db {
	//	return Db{}
	//}
	//
	//type NameProvider struct {
	//	db Db
	//}
	//
	//func (db Db) GetName() string {
	//	return "Elton John"
	//}
	//
	//func NewNameProvider(db Db) NameProvider {
	//	return NameProvider{db: db}
	//}
	//
	//func (dp NameProvider) GetFullName(salutation string) string {
	//	return salutation + " " + dp.db.GetName()
	//}

	cont := NewRuntimeContainer()
	cont.AddNewMethod("db", NewDb)
	cont.AddNewMethod("nameProvider", NewNameProvider, "db")

	nameProvider := cont.Get("nameProvider", true).(NameProvider)
	fmt.Println(nameProvider.GetFullName("Sir"))

	// Output:
	// Sir Elton John
}
