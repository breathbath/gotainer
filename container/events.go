package container

type dependencyNotifier func(observer interface{}, dependency interface{})

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

func (this *EventsContainer) registerDependencyEvent(eventName, dependencyName string) {
	if this.dependencyEvents[eventName] == nil {
		this.dependencyEvents[eventName] = []string{}
	}
	this.dependencyEvents[eventName] = append(this.dependencyEvents[eventName], dependencyName)
}

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
