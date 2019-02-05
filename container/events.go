package container

//serviceNotificationCallback is a function that receives Observer as a Service interested in a dependency
//received in the second argument so you can call it as Observer.SetSomeDependency(dependency)
type serviceNotificationCallback func(serviceInterestedInDependency interface{}, dependency interface{}) error

//EventsContainer contains all Observer, events and Config declarations, that you might
//add in your container
type EventsContainer struct {
	dependencyEvents             map[string][]string
	serviceNotificationCallbacks map[string]map[string]serviceNotificationCallback
}

//NewEventsContainer EventsContainer Constr
func NewEventsContainer() *EventsContainer {
	return &EventsContainer{
		dependencyEvents:             map[string][]string{},
		serviceNotificationCallbacks: map[string]map[string]serviceNotificationCallback{},
	}
}

//registerDependencyEvent triggers an Event about adding a concrete dependency to the container
func (ec *EventsContainer) registerDependencyEvent(eventName, dependencyName string) {
	ec.initEventCollection(eventName)
	ec.dependencyEvents[eventName] = append(ec.dependencyEvents[eventName], dependencyName)
}

//addDependencyObserver adds the Service (Observer) which will receive Config added by known events
func (ec *EventsContainer) addDependencyObserver(
	eventName,
	serviceId string,
	callbackToProvideDependencyToService interface{},
) error {
	if ec.serviceNotificationCallbacks[serviceId] == nil {
		ec.serviceNotificationCallbacks[serviceId] = map[string]serviceNotificationCallback{}
	}
	notifCallack, err := wrapCallbackToProvideDependencyToServiceIntoServiceNotificationCallback(
		callbackToProvideDependencyToService,
		eventName,
		serviceId,
	)
	if err != nil {
		return err
	}

	ec.serviceNotificationCallbacks[serviceId][eventName] = notifCallack
	return nil
}

//collectDependencyEventsForService we call Observer methods with all the Config that it's interested in
func (ec *EventsContainer) collectDependencyEventsForService(
	c Container,
	serviceId string,
	serviceInstance interface{},
) error {
	dependencyEventSubscribedServices, eventObserverFound := ec.serviceNotificationCallbacks[serviceId]
	if !eventObserverFound {
		return nil
	}

	errs := []error{}
	for eventName, serviceNotificationCallback := range dependencyEventSubscribedServices {
		dependencies, eventFound := ec.dependencyEvents[eventName]
		if !eventFound {
			continue
		}

		for _, dependencyName := range dependencies {
			dependency := c.Get(dependencyName, true)
			err := serviceNotificationCallback(serviceInstance, dependency)
			if err != nil {
				errs = append(errs, err)
			}
		}
	}

	return mergeErrors(errs)
}

//merge helps to accumulate Event collections when we try to merge containers
func (ec *EventsContainer) merge(ecToCopy EventsContainer) error {
	for ecKey, events := range ecToCopy.dependencyEvents {
		for _, dependencyName := range events {
			ec.registerDependencyEvent(ecKey, dependencyName)
		}
	}

	errs := []error{}
	for observerId, dependencyNotifiers := range ecToCopy.serviceNotificationCallbacks {
		for eventName, dependencyNotifier := range dependencyNotifiers {
			err := ec.addDependencyObserver(eventName, observerId, dependencyNotifier)
			if err != nil {
				errs = append(errs, err)
			}
		}
	}

	return mergeErrors(errs)
}

func (ec *EventsContainer) initEventCollection(eventName string) {
	if ec.dependencyEvents[eventName] == nil {
		ec.dependencyEvents[eventName] = []string{}
	}
}
