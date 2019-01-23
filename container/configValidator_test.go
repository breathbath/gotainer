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

func TestNoConstructorNorNewMethodProvidedSecureMode(t *testing.T) {
	node := Node{
		ID: "someService",
		GarbageFunc: func(service interface{}) error {
			return nil
		},
	}
	assertWrongNodeDeclarationSecure(
		node,
		t,
		"A new or constructor function are expected but none was declared [check '%s' service]",
		"someService",
	)
}

func TestNoConstructorNorNewMethodNorIdProvidedSecureMode(t *testing.T) {
	node := Node{
		ID: "",
	}
	assertWrongNodeDeclarationSecure(
		node,
		t,
		"A new or constructor function are expected but none was declared see '%s'",
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
		"A new or constructor function are expected but none was declared [check 'serviceNamesWithoutNewFunc' service];\nServices list should be defined with a non empty new func, see '%s'",
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
		Ev: Event{Name: "add_stats_provider"},
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
			Ev: Event{Name: "add_stats_provider", Service: "Some unknown service"},
		},
		t,
		"Unknown service declaration 'Some unknown service' in 'event add_stats_provider'",
	)
}

func TestConfigValidationIsTriggeredWithContainerBuilder(t *testing.T) {
	configTree := getMockedConfigTree()
	node := Node{
		ID:      "",
		NewFunc: mocks.NewConfig,
	}
	configTree = append(configTree, node)

	_, err := RuntimeContainerBuilder{}.BuildContainerFromConfigSecure(configTree)

	assertErrorText(
		fmt.Sprintf("The new function should be provided with a service id, see '%s'", node),
		err,
		t,
	)
}

//the implementation might trigger validation for each tree separately but this is not expected because
//a subtree might reference a service in another tree which is perfectly valid
func TestConfigValidationForMergedConfigTreesInBuilder(t *testing.T) {
	configTree1 := getMockedConfigTree()
	node1 := Node{
		Ev:      Event{
			Name: "event_me",
			Service: "some_service_from_config_tree_2",
		},
	}
	node2 := Node{
		Ob:      Observer{
			Event: "event_me",
			Name: "book_creator",
			Callback: func(bookCreator mocks.BookCreator, someStr string) {},
		},
	}
	configTree1 = append(configTree1, node1, node2)

	_, err := RuntimeContainerBuilder{}.BuildContainerFromConfigSecure(configTree1)

	assertErrorText(
		"Unknown service declaration 'some_service_from_config_tree_2' in 'event event_me'",
		err,
		t,
	)

	configTree2 := Tree{
		Node{
			ID: "some_service_from_config_tree_2",
			Constr: func(c Container) (i interface{}, e error) {
				return "some service", nil
			},
		},
	}

	_, err = RuntimeContainerBuilder{}.BuildContainerFromConfigSecure(configTree1, configTree2)
	assertNoError(err, t)
}

func TestEventWithoutCorrespondingObserver(t *testing.T) {
	configTree := getMockedConfigTree()
	node := Node{
		Ev:      Event{
			Name: "some_book_event",
			Service: "book_storage",
		},
	}
	configTree = append(configTree, node)

	assertWrongNodeDeclaration(
		node,
		t,
		"No observer is declared for the event 'some_book_event'",
	)
}

func assertWrongNodeDeclarationSecure(node Node, t *testing.T, expectedErrorFormat string, context ...interface{}) {
	configTree := buildConfigTree(node)

	expectedErrText := buildExpectedErrorText(expectedErrorFormat, context)
	err := ValidateConfigSecure(configTree)

	assertErrorText(expectedErrText, err, t)
}

func assertWrongNodeDeclaration(node Node, t *testing.T, expectedErrorFormat string, context ...interface{}) {
	expectedErrText := buildExpectedErrorText(expectedErrorFormat, context)
	defer ExpectPanic(t, expectedErrText)

	configTree := buildConfigTree(node)

	ValidateConfig(configTree)
}

func buildConfigTree(node Node) Tree {
	configTree := getMockedConfigTree()
	configTree = append(configTree, node)

	return configTree
}

func buildExpectedErrorText(expectedErrorFormat string, context []interface{}) string {
	var errorText string
	if len(context) > 0 {
		errorText = fmt.Sprintf(expectedErrorFormat, context...)
	} else {
		errorText = expectedErrorFormat
	}

	return errorText
}
