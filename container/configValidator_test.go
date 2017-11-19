package container

import (
	"testing"
	"github.com/breathbath/gotainer/container/mocks"
)

func TestNewFunctionIsNotFunction(t *testing.T) {
	assertWrongNodeDeclaration(
		"A function is expected rather than 'string' in the service declaration 'wrongNewAsString'",
		"wrongNewAsString",
		Node{NewFunc: "abc"},
		t,
	)
}

func TestNewFunctionReturnsNoValues(t *testing.T) {
	assertWrongNodeDeclaration(
		"Constr function should return 1 or 2 values, but 0 values are returned [check 'wrongNewWrongReturnValues1' service]",
		"wrongNewWrongReturnValues1",
		Node{NewFunc: func() {}},
		t,
	)
}

func TestNewFunctionReturnsNoErrors(t *testing.T) {
	assertWrongNodeDeclaration(
		"Constr function with 2 returned values should return at least one error interface [check 'wrongNewWrongReturnValues2' service]",
		"wrongNewWrongReturnValues2",
		Node{NewFunc: func() (string, int) {
			return "", 0
		}},
		t,
	)
}

func TestNewFunctionHasLessArgumentsThanServiceNamesCount(t *testing.T) {
	assertWrongNodeDeclaration(
		"The function requires 1 arguments, but 2 arguments are provided in the service declaration 'wrongArgumentsCount'",
		"wrongArgumentsCount",
		Node{
			NewFunc: func(a string) string {
				return a
			},
			ServiceNames: Services{"a", "b"},
		},
		t,
	)
}

func TestMoreDefinitionsForNewFunction(t *testing.T) {
	assertWrongNodeDeclaration(
		`Unexpected constructor declaration, see service 'moreDefinitionsForNewFunction';
Unexpected event declaration, see service 'moreDefinitionsForNewFunction'`,
		"moreDefinitionsForNewFunction",
		Node{
			NewFunc: mocks.NewConfig,
			Ev:      Event{Name: "someEvent"},
			Constr: func(c Container) (interface{}, error) {
				return nil, nil
			},
		},
		t,
	)
}

func TestMoreDeclarationsForConstrFunction(t *testing.T) {
	assertWrongNodeDeclaration(
		"Unexpected observer declaration, see service 'moreDefinitionsForConstrFunction'",
		"moreDefinitionsForConstrFunction",
		Node{
			NewFunc: mocks.NewConfig,
			Ob:      Observer{Name: "someEvent"},
		},
		t,
	)
}

func TestServiceNamesProvidedWithoutNewFunc(t *testing.T) {
	assertWrongNodeDeclaration(
		"Services list should be defined with a non empty new func, see service 'serviceNamesWithoutNewFunc'",
		"serviceNamesWithoutNewFunc",
		Node{
			ServiceNames: []string{"someService"},
		},
		t,
	)
}

func TestObserverRequiredFieldsNotProvided(t *testing.T) {
	assertWrongNodeDeclaration(
		`Observer name is required, see service 'observerRequiredFields';
A function is expected rather than 'string' in the service declaration 'observerRequiredFields'`,
		"observerRequiredFields",
		Node{
			Ob: Observer{Event: "someEv", Name: "", Callback: "lsls"},
		},
		t,
	)

	assertWrongNodeDeclaration(
		`Observer event is required, see service 'observerRequiredFields'`,
		"observerRequiredFields",
		Node{
			Ob: Observer{Event: "", Name: "someName", Callback: func(sg *mocks.StatisticsGateway, sp mocks.StatisticsProvider){}},
		},
		t,
	)

	assertWrongNodeDeclaration(
		`Observer callback is required, see service 'observerRequiredFields';
A function is expected rather than 'invalid' in the service declaration 'observerRequiredFields'`,
		"observerRequiredFields",
		Node{
			Ob: Observer{Event: "someEv", Name: "someName", Callback: nil},
		},
		t,
	)
}

func TestMoreDeclarationsForObserver(t *testing.T) {
	assertWrongNodeDeclaration(
		"Unexpected event declaration, see service 'moreDefinitionsForObserver'",
		"moreDefinitionsForObserver",
		Node{
			Ob: Observer{Event: "someEv", Name: "someName", Callback: func(sg *mocks.StatisticsGateway, sp mocks.StatisticsProvider){}},
			Ev:      Event{Name:"someEvent"},
		},
		t,
	)
}

func TestEventRequiredFieldsNotProvided(t *testing.T) {
	assertWrongNodeDeclaration(
		`Event name is required, see service 'eventRequiredFields'`,
		"eventRequiredFields",
		Node{
			Ev: Event{Service: "config"},
		},
		t,
	)

	assertWrongNodeDeclaration(
		`Event service is required, see service 'eventRequiredFields'`,
		"eventRequiredFields",
		Node{
			Ev: Event{Name: "someEvent"},
		},
		t,
	)
}

func TestUnknownEventServiceIsProvided(t *testing.T) {
	assertWrongNodeDeclaration(
		`Unknown service declaration 'Some unknown service' in 'event someEvent'`,
		"unknownEventService",
		Node{
			Ev: Event{Name: "someEvent", Service:"Some unknown service"},
		},
		t,
	)
}

func assertWrongNodeDeclaration(expectedPanicMessage, serviceId string, node Node, t *testing.T) {
	configTree := getMockedConfigTree()
	configTree[serviceId] = node

	defer ExpectPanic(expectedPanicMessage, t)

	ValidateConfig(configTree)
}
