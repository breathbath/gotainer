package container

import (
	"fmt"
	"github.com/breathbath/gotainer/container/mocks"
	"testing"
)

func TestDefaultIsValid(t *testing.T) {
	configTree := getMockedConfigTree()

	ValidateConfig(configTree)
}

func TestNewFunctionIsNotFunction(t *testing.T) {
	node := Node{NewFunc: "abc", ID: "wrongNewAsString"}
	assertWrongNodeDeclaration(
		node,
		t,
		"A function is expected rather than 'string' [check '%s' service]",
		node,
	)
}

func TestNewFunctionReturnsNoValues(t *testing.T) {
	assertWrongNodeDeclaration(
		Node{
			NewFunc: func() {},
			ID:      "wrongNewWrongReturnValues1",
		},
		t,
		"Constr function should return 1 or 2 values, but 0 values are returned [check 'wrongNewWrongReturnValues1' service]",
	)
}

func TestNewFunctionReturnsNoErrors(t *testing.T) {
	assertWrongNodeDeclaration(
		Node{
			ID: "wrongNewWrongReturnValues2",
			NewFunc: func() (string, int) {
				return "", 0
			},
		},
		t,
		"Constr function with 2 returned values should return at least one error interface [check 'wrongNewWrongReturnValues2' service]",
	)
}

func TestNewFunctionHasLessArgumentsThanServiceNamesCount(t *testing.T) {
	node := Node{
		ID: "wrongArgumentsCount",
		NewFunc: func(a string) string {
			return a
		},
		ServiceNames: Services{"a", "b"},
	}
	assertWrongNodeDeclaration(
		node,
		t,
		"The function requires 1 arguments, but 2 arguments are provided [check '%s' service]",
		node,
	)
}

func TestNewFunctionMissingId(t *testing.T) {
	node := Node{
		ID:      "",
		NewFunc: mocks.NewConfig,
	}

	assertWrongNodeDeclaration(
		node,
		t,
		"The new function should be provided with a service id, see '%s'",
		node,
	)
}

func TestMoreDefinitionsForNewFunction(t *testing.T) {
	node := Node{
		ID:      "moreDefinitionsForNewFunction",
		NewFunc: mocks.NewConfig,
		Ev:      Event{Name: "someEvent"},
		Constr: func(c Container) (interface{}, error) {
			return nil, nil
		},
	}
	assertWrongNodeDeclaration(
		node,
		t,
		"Unexpected constructor declaration, see '%s';\nUnexpected event declaration, see '%s'",
		node,
		node,
	)
}

func TestMoreDeclarationsForConstrFunction(t *testing.T) {
	node := Node{
		ID:      "moreDefinitionsForConstrFunction",
		NewFunc: mocks.NewConfig,
		Ob:      Observer{Name: "someEvent"},
	}
	assertWrongNodeDeclaration(
		node,
		t,
		"Unexpected observer declaration, see '%s'",
		node,
	)
}

func TestNoIdForConstrFunction(t *testing.T) {
	node := Node{
		ID: "",
		Constr: func(c Container) (interface{}, error) {
			return "", nil
		},
	}

	assertWrongNodeDeclaration(
		node,
		t,
		"The constructor function should be provided with a non empty service id, see '%s'",
		node,
	)
}

func TestServiceNamesProvidedWithoutNewFunc(t *testing.T) {
	node := Node{
		ID:           "serviceNamesWithoutNewFunc",
		ServiceNames: []string{"someService"},
	}
	assertWrongNodeDeclaration(
		node,
		t,
		"Services list should be defined with a non empty new func, see '%s'",
		node,
	)
}

func TestObserverRequiredFieldsNotProvided(t *testing.T) {
	node := Node{
		Ob: Observer{Event: "someEv", Name: "", Callback: "lsls"},
	}
	assertWrongNodeDeclaration(
		node,
		t,
		"Observer name is required [check '%s' service];\nA function is expected rather than 'string' [check '%s' service]",
		node,
		node,
	)

	node = Node{
		Ob: Observer{Event: "", Name: "someName", Callback: func(sg *mocks.StatisticsGateway, sp mocks.StatisticsProvider) {}},
	}
	assertWrongNodeDeclaration(
		node,
		t,
		"Observer event is required [check '%s' service]",
		node,
	)

	node = Node{
		Ob: Observer{Event: "someEv", Name: "someName", Callback: nil},
	}
	assertWrongNodeDeclaration(
		node,
		t,
		"Observer callback is required [check '%s' service];\nA function is expected rather than 'invalid' [check '%s' service]",
		node,
		node,
	)
}

func TestMoreDeclarationsForObserver(t *testing.T) {
	node := Node{
		Ob: Observer{Event: "someEv", Name: "someName", Callback: func(sg *mocks.StatisticsGateway, sp mocks.StatisticsProvider) {}},
		Ev: Event{Name: "someEvent"},
	}
	assertWrongNodeDeclaration(
		node,
		t,
		"Unexpected event declaration, see '%s'",
		node,
	)
}

func TestEventRequiredFieldsNotProvided(t *testing.T) {
	node := Node{
		Ev: Event{Service: "config"},
	}
	assertWrongNodeDeclaration(
		node,
		t,
		"Event name is required, see '%s'",
		node,
	)

	node = Node{
		Ev: Event{Name: "someEvent"},
	}
	assertWrongNodeDeclaration(
		node,
		t,
		"Event service is required, see '%s'",
		node,
	)
}

func TestUnknownEventServiceIsProvided(t *testing.T) {
	assertWrongNodeDeclaration(
		Node{
			ID: "unknownEventService",
			Ev: Event{Name: "someEvent", Service: "Some unknown service"},
		},
		t,
		"Unknown service declaration 'Some unknown service' in 'event someEvent'",
	)
}

func assertWrongNodeDeclaration(node Node, t *testing.T, expectedErrorFormat string, context ...interface{}) {
	configTree := getMockedConfigTree()
	configTree = append(configTree, node)
	var errorText string
	if len(context) > 0 {
		errorText = fmt.Sprintf(expectedErrorFormat, context...)
	} else {
		errorText = expectedErrorFormat
	}
	defer ExpectPanic(t, errorText)

	ValidateConfig(configTree)
}
