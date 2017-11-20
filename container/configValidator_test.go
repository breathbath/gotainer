package container

import (
	"testing"
	"github.com/breathbath/gotainer/container/mocks"
	"fmt"
)

func TestDefaultIsValid(t *testing.T) {
	configTree := getMockedConfigTree()

	ValidateConfig(configTree)
}

func TestNewFunctionIsNotFunction(t *testing.T) {
	assertWrongNodeDeclaration(
		"A function is expected rather than 'string' in the service declaration 'wrongNewAsString'",
		Node{NewFunc: "abc", Id: "wrongNewAsString"},
		t,
	)
}

func TestNewFunctionReturnsNoValues(t *testing.T) {
	assertWrongNodeDeclaration(
		"Constr function should return 1 or 2 values, but 0 values are returned [check 'wrongNewWrongReturnValues1' service]",
		Node{
			NewFunc: func() {},
			Id: "wrongNewWrongReturnValues1",
			},
		t,
	)
}

func TestNewFunctionReturnsNoErrors(t *testing.T) {
	assertWrongNodeDeclaration(
		"Constr function with 2 returned values should return at least one error interface [check 'wrongNewWrongReturnValues2' service]",
		Node{
			Id: "wrongNewWrongReturnValues2",
			NewFunc: func() (string, int) {
				return "", 0
			},
		},
		t,
	)
}

func TestNewFunctionHasLessArgumentsThanServiceNamesCount(t *testing.T) {
	assertWrongNodeDeclaration(
		"The function requires 1 arguments, but 2 arguments are provided in the service declaration 'wrongArgumentsCount'",
		Node{
			Id: "wrongArgumentsCount",
			NewFunc: func(a string) string {
				return a
			},
			ServiceNames: Services{"a", "b"},
		},
		t,
	)
}

func TestNewFunctionMissingId(t *testing.T) {
	node := Node{
		Id: "",
		NewFunc: mocks.NewConfig,
	}

	nodeText := fmt.Sprint(node)
	assertWrongNodeDeclaration(
		"The new function should be provided with a service id, see '" + nodeText + "'",
		node,
		t,
	)
}

func TestMoreDefinitionsForNewFunction(t *testing.T) {
	assertWrongNodeDeclaration(
		`Unexpected constructor declaration, see service 'moreDefinitionsForNewFunction';
Unexpected event declaration, see service 'moreDefinitionsForNewFunction'`,
		Node{
			Id:      "moreDefinitionsForNewFunction",
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
		Node{
			Id:      "moreDefinitionsForConstrFunction",
			NewFunc: mocks.NewConfig,
			Ob:      Observer{Name: "someEvent"},
		},
		t,
	)
}

func TestNoIdForConstrFunction(t *testing.T) {
	node := 		Node{
		Id:      "",
		Constr: func(c Container) (interface{}, error) {
			return "", nil
		},
	}

	assertWrongNodeDeclaration(
		"The constructor function should be provided with a non empty service id, see '" + fmt.Sprint(node) + "'",
		node,
		t,
	)
}

func TestServiceNamesProvidedWithoutNewFunc(t *testing.T) {
	assertWrongNodeDeclaration(
		"Services list should be defined with a non empty new func, see service 'serviceNamesWithoutNewFunc'",
		Node{
			Id:           "serviceNamesWithoutNewFunc",
			ServiceNames: []string{"someService"},
		},
		t,
	)
}

func TestObserverRequiredFieldsNotProvided(t *testing.T) {
	assertWrongNodeDeclaration(
		`Observer name is required, see observer '';
A function is expected rather than 'string' in the observer declaration ''`,
		Node{
			Ob: Observer{Event: "someEv", Name: "", Callback: "lsls"},
		},
		t,
	)

	assertWrongNodeDeclaration(
		`Observer event is required, see observer 'someName'`,
		Node{
			Ob: Observer{Event: "", Name: "someName", Callback: func(sg *mocks.StatisticsGateway, sp mocks.StatisticsProvider) {}},
		},
		t,
	)

	assertWrongNodeDeclaration(
		`Observer callback is required, see service 'observerRequiredFields';
A function is expected rather than 'invalid' in the observer declaration 'observerRequiredFields'`,
		Node{
			Ob: Observer{Event: "someEv", Name: "someName", Callback: nil},
		},
		t,
	)
}

func TestMoreDeclarationsForObserver(t *testing.T) {
	assertWrongNodeDeclaration(
		"Unexpected event declaration, see observer 'someName'",
		Node{
			Ob: Observer{Event: "someEv", Name: "someName", Callback: func(sg *mocks.StatisticsGateway, sp mocks.StatisticsProvider) {}},
			Ev: Event{Name: "someEvent"},
		},
		t,
	)
}

func TestEventRequiredFieldsNotProvided(t *testing.T) {
	assertWrongNodeDeclaration(
		`Event name is required, see event 'eventRequiredFields'`,
		Node{
			Ev: Event{Service: "config"},
		},
		t,
	)

	assertWrongNodeDeclaration(
		`Event service is required, see event 'someEvent'`,
		Node{
			Ev: Event{Name: "someEvent"},
		},
		t,
	)
}

func TestUnknownEventServiceIsProvided(t *testing.T) {
	assertWrongNodeDeclaration(
		`Unknown service declaration 'Some unknown service' in 'event someEvent'`,
		Node{
			Id: "unknownEventService",
			Ev: Event{Name: "someEvent", Service: "Some unknown service"},
		},
		t,
	)
}

func assertWrongNodeDeclaration(expectedPanicMessage string, node Node, t *testing.T) {
	configTree := getMockedConfigTree()
	configTree = append(configTree, node)

	defer ExpectPanic(expectedPanicMessage, t)

	ValidateConfig(configTree)
}
