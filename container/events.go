package container

//serviceNotificationCallback is a function that receives observer as a service interested in a dependency
//received in the second argument so you can call it as observer.SetSomeDependency(dependency)
type serviceNotificationCallback func(serviceInterestedInDependency interface{}, dependency interface{})

//EventsContainer contains all observer, events and dependencies declarations, that you might
//add in your container
type EventsContainer struct {
	dependencyEvents             map[string][]string
	serviceNotificationCallbacks map[string]map[string]serviceNotificationCallback
}

//NewEventsContainer EventsContainer constructor
func NewEventsContainer() *EventsContainer {
	return &EventsContainer{
		dependencyEvents:             map[string][]string{},
		serviceNotificationCallbacks: map[string]map[string]serviceNotificationCallback{},
	}
}

//registerDependencyEvent triggers an event about adding a concrete dependency to the container
func (ec *EventsContainer) registerDependencyEvent(eventName, dependencyName string) {
	ec.initEventCollection(eventName)
	ec.dependencyEvents[eventName] = append(ec.dependencyEvents[eventName], dependencyName)
}

//addDependencyObserver adds the service (observer) which will receive dependencies added by known events
func (ec *EventsContainer) addDependencyObserver(eventName, serviceId string, callbackToProvideDependencyToService interface{}) {
	if ec.serviceNotificationCallbacks[serviceId] == nil {
		ec.serviceNotificationCallbacks[serviceId] = map[string]serviceNotificationCallback{}
	}
	ec.serviceNotificationCallbacks[serviceId][eventName] = wrapCallbackToProvideDependencyToServiceIntoServiceNotificationCallback(
		callbackToProvideDependencyToService,
		eventName,
		serviceId,
	)
}

//collectDependencyEventsForService we call observer methods with all the dependencies that it's interested in
func (ec *EventsContainer) collectDependencyEventsForService(c Container, serviceId string, serviceInstance interface{}) {
	dependencyEventSubscribedServices, eventObserverFound := ec.serviceNotificationCallbacks[serviceId]
	if !eventObserverFound {
		return
	}

	for eventName, serviceNotificationCallback := range dependencyEventSubscribedServices {
		dependencies, eventFound := ec.dependencyEvents[eventName]
		if !eventFound {
			continue
		}

		for _, dependencyName := range dependencies {
			dependency := c.Get(dependencyName, true)
			serviceNotificationCallback(serviceInstance, dependency)
		}
	}
}

//merge helps to accumulate event collections when we try to merge containers
func (ec *EventsContainer) merge(ecToCopy EventsContainer) {
	for ecKey, events := range ecToCopy.dependencyEvents {
		for _, dependencyName := range events {
			ec.registerDependencyEvent(ecKey, dependencyName)
		}
	}

	for observerId, dependencyNotifiers := range ecToCopy.serviceNotificationCallbacks {
		for eventName, dependencyNotifier := range dependencyNotifiers {
			ec.addDependencyObserver(eventName, observerId, dependencyNotifier)
		}
	}
}

func (ec *EventsContainer) initEventCollection(eventName string) {
	if ec.dependencyEvents[eventName] == nil {
		ec.dependencyEvents[eventName] = []string{}
	}
}
