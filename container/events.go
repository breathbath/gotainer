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

func NewEventsContainer() *EventsContainer {
	return &EventsContainer{
		dependencyEvents:    map[string][]string{},
		dependencyObservers: map[string]map[string]dependencyNotifier{},
	}
}

//registerDependencyEvent triggers an event about adding a concrete dependency to the container
func (this *EventsContainer) registerDependencyEvent(eventName, dependencyName string) {
	if this.dependencyEvents[eventName] == nil {
		this.dependencyEvents[eventName] = []string{}
	}
	this.dependencyEvents[eventName] = append(this.dependencyEvents[eventName], dependencyName)
}

//addDependencyObserver adds the service (observer) which will receive dependencies added by known events
func (this *EventsContainer) addDependencyObserver(eventName, observerId string, observerResolver interface{}) {
	if this.dependencyObservers[observerId] == nil {
		this.dependencyObservers[observerId] = map[string]dependencyNotifier{}
	}
	this.dependencyObservers[observerId][eventName] = convertCustomObserverResolverToDependencyNotifier(
		observerResolver,
		eventName,
		observerId,
	)
}

//notifyObserverAboutDependency we call observer methods with all the dependencies that it's interested in
func (this *EventsContainer) notifyObserverAboutDependency(c RuntimeContainer, observerId string, observer interface{}) {
	eventObservers, eventObserverFound := this.dependencyObservers[observerId]
	if !eventObserverFound {
		return
	}

	for eventName, dependencyObserver := range eventObservers {
		dependencies, eventFound := this.dependencyEvents[eventName]
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
func (this *EventsContainer) merge(ec EventsContainer) {
	for ecKey, events := range ec.dependencyEvents {
		for _, dependencyName := range events {
			this.dependencyEvents[ecKey] = append(this.dependencyEvents[ecKey], dependencyName)
		}
	}

	for observerId, dependencyNotifiers := range ec.dependencyObservers {
		for eventName, dependencyNotifier := range dependencyNotifiers {
			this.dependencyObservers[observerId][eventName] = dependencyNotifier
		}
	}
}
