# Dependencies injection tool for go driven projects

[![Travis Status for breathbath/gotainer](https://api.travis-ci.org/breathbath/gotainer.svg?branch=master&label=linux+build)](https://travis-ci.org/breathbath/gotainer)
[![godoc for breathbath/gotainer](https://godoc.org/github.com/nathany/looper?status.svg)](http://godoc.org/github.com/breathbath/gotainer/container)
[![goreportcard for breathbath/gotainer](https://goreportcard.com/badge/github.com/breathbath/gotainer)](https://goreportcard.com/report/breathbath/gotainer)
[![codecov for breathbath/gotainer](https://codecov.io/gh/breathbath/gotainer/branch/master/graph/badge.svg)](https://codecov.io/gh/breathbath/gotainer)

This library helps to manage dependencies in your project by providing a centralised logic for initialising services.

# Use cases

The library covers the following use cases:

## Reusable services with dependencies

Imagine we have 2 services dependant on an another one, which should be reused in both cases. The common service
should be initialised only once as it has an internal state (e.g., db connection).

        //simple service with no dependencies
        type ServiceA struct {}

        //simple service with a dependency
        type ServiceB struct {
            serviceA ServiceA
        }

        //more complex service depending on 2 others, which are also dependant
        type ServiceC struct {
           serviceA ServiceA
           serviceB ServiceB
        }

        //services declaration
        container.AddNewMethod("service_a", package_a.New)
        container.AddNewMethod("service_b", package_b.New, "service_a")
        container.AddNewMethod("service_c", package_c.New, "service_a", "service_b")

        //service fetching, here you can enjoy the fully typed service declaration
        var serviceC ServiceC
        container.scan("service_c", &serviceC)

## Services with a complex initialisation

        type ServiceX struct {...}
        func (sx ServiceX) AddService(sc ServiceC){}
        func (sx ServiceX) EnableLogging(){}
        func (sx ServiceX) RegisterInMonitoringList(monitoringList []MonitoringItem){}

        //we can do a complex ServiceX initialisation, this code will be executed once and all services using
        //"service_x" will have a fully initialised version of it
        container.AddConstructor("service_x",  func(c container.Container) (interface{}, error){
            serviceX := NewServiceX

            var serviceC ServiceC
            c.Scan("service_c", &serviceC)
            serviceX.AddService(serviceC)

            serviceX.EnableLogging()

            var monitoringList []MonitoringItem
            c.Scan("monitoring_list", &monitoringList)
            serviceX.RegisterInMonitoringList(monitoringList)

            return serviceX, nil
        })

## Cached and reusable results of a method call or parameters as dependencies

        func CountItems() int64 {...}

        //in this case CountItems can be an expensive operation that should be executed once
        container.AddConstructor("items_count", func(c container.Container) (interface{}, error) {
            return CountItems(), nil
        })

        //we declare a simple string config option as a container service
        container.AddConstructor("static_url", func(c container.Container) (interface{}, error) {
            var config Config
            c.Scan("config", &config)

            return config.GetValue("static_url"), nil
        })

## Explicitly non cached services

Sometimes we want to recreate a service every time we fetch it:

        var serviceA ServiceA
        container.ScanNonCached("service_a", &serviceA)

## Services chain

        container.AddConstructor("chained_services", func(c container.Container) (interface{}, error) {
               var initialService StartingPoint
               c.Scan("startingPoint", &initialService)
               return initialService.GetA().GetB().GetC(), nil
        })

## Anonymous constructors

        //you actually don't need "new" methods for your services
        container.AddConstructor("service_a", func(c container.Container) (interface{}, error) {
               return ServiceA{}
        })

## Dependency events

In some cases your service should get certain dependencies every time when they are added to the container. This logic
helps to to avoid multiple calls of the same method on your service and also detach new dependant services registration
from your main service. Consider following example:

        type MonitoringProvider interface{
            GetMonitoringEvent() (eventName string, count int64)
        }

        type TotalMonitoringProvider struct{...}
        func(tmp TotalMonitoringProvider) GetMonitoringEvent{
            return "total_count", 100
        }

        type ErrorCountProvider struct{...}
        func(tmp ErrorCountProvider) GetMonitoringEvent{
            return "errors_count", 10
        }

        type MonitoringGateway struct{...}
        func (mg MonitoringGateway) AddMonitoringProvider (mp MonitoringProvider)...


In this case we expect new implementations of MonitoringProvider will be added in future. An obvious solution would be:

        func BuildMonitoringGateway(tmp TotalMonitoringProvider, ecp ErrorCountProvider) MonitoringGateway {
            mg := MonitoringGateway{}
            mg.AddMonitoringProvider(tmp)
            mg.AddMonitoringProvider(ecp)
        }

This approach has following problems:

1. With every new implementation of MonitoringProvider you should modify the BuildMonitoringGateway, so this code is not closed to
modification.

2. The amount of arguments of BuildMonitoringGateway will grow, so this function becomes unreadable

3. You should create every new instance of MonitoringProvider somewhere which will probably lead to code duplication, if those
require other services, the amount of boilerplate code will explode

With the Gotainer you can solve this problem with the following code:

        container.AddDependencyObserver("monitoring_provided_added", "monitoring_gateway", func(mg MonitoringGateway, mp MonitoringProvider){
            sg.AddMonitoringProvider(sp)
        })

        container.AddNewMethod("total_monitoring_provider", NewTotalMonitoringProvider, "service_a", "service_b")
        container.RegisterDependencyEvent("monitoring_provided_added", "total_monitoring_provider")

        container.AddNewMethod("error_count_provider", NewErrorCountProvider)
        container.RegisterDependencyEvent("monitoring_provided_added", "error_count_provider")

This has following advantages:

1. MonitoringGateway is completely decoupled from adding new implementations of MonitoringProvider

2. No complex initialisation function for MonitoringGateway is needed

3. Concrete implementations of MonitoringProvider are created once without any repetition as this logic is already encapsulated in the Gotainer.

4. You might have the container declaration for your MonitoringGateway in one core library package and different implementations of
MonitoringProvider in other packages, so you are able to plug them in individually in every application with no need to change the
core package.