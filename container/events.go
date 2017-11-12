package container

//dependencyNotifier is a function that receives observer as a service interested in a dependency
//received in the second argument so you can call it as observer.SetSomeDependency(dependency)
type dependencyNotifier func(observer interface{}, dependency interface{})

//EventsContainer contains all observer, events and dependencies declarations, that you might
//add in your container
type EventsContainer struct {
	dependencyEvents    map[string][]string
	dependencyObservers map[string]map[string]dependencyNotifier
}

//NewEventsContainer EventsContainer constructor
func NewEventsContainer() *EventsContainer {
	return &EventsContainer{
		dependencyEvents:    map[string][]string{},
		dependencyObservers: map[string]map[string]dependencyNotifier{},
	}
}

//registerDependencyEvent triggers an event about adding a concrete dependency to the container
func (ec *EventsContainer) registerDependencyEvent(eventName, dependencyName string) {
	if ec.dependencyEvents[eventName] == nil {
		ec.dependencyEvents[eventName] = []string{}
	}
	ec.dependencyEvents[eventName] = append(ec.dependencyEvents[eventName], dependencyName)
}

//addDependencyObserver adds the service (observer) which will receive dependencies added by known events
func (ec *EventsContainer) addDependencyObserver(eventName, observerId string, observerResolver interface{}) {
	if ec.dependencyObservers[observerId] == nil {
		ec.dependencyObservers[observerId] = map[string]dependencyNotifier{}
	}
	ec.dependencyObservers[observerId][eventName] = convertCustomObserverResolverToDependencyNotifier(
		observerResolver,
		eventName,
		observerId,
	)
}

//notifyObserverAboutDependency we call observer methods with all the dependencies that it's interested in
func (ec *EventsContainer) notifyObserverAboutDependency(c RuntimeContainer, observerId string, observer interface{}) {
	eventObservers, eventObserverFound := ec.dependencyObservers[observerId]
	if !eventObserverFound {
		return
	}

	for eventName, dependencyObserver := range eventObservers {
		dependencies, eventFound := ec.dependencyEvents[eventName]
		if !eventFound {
			continue
		}

		for _, dependencyName := range dependencies {
			dependency := c.Get(dependencyName, true)
			dependencyObserver(observer, dependency)
		}
	}
}

//merge helps to accumulate event collections when we try to merge containers
func (ec *EventsContainer) merge(ecToCopy EventsContainer) {
	for ecKey, events := range ecToCopy.dependencyEvents {
		for _, dependencyName := range events {
			ec.dependencyEvents[ecKey] = append(ec.dependencyEvents[ecKey], dependencyName)
		}
	}

	for observerId, dependencyNotifiers := range ecToCopy.dependencyObservers {
		for eventName, dependencyNotifier := range dependencyNotifiers {
			ec.dependencyObservers[observerId][eventName] = dependencyNotifier
		}
	}
}
