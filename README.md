# Dependency container for go driven projects

[![Travis Status for breathbath/gotainer](https://api.travis-ci.org/breathbath/gotainer.svg?branch=master&label=linux+build)](https://travis-ci.org/breathbath/gotainer)
[![godoc for breathbath/gotainer](https://godoc.org/github.com/nathany/looper?status.svg)](http://godoc.org/github.com/breathbath/gotainer/container)
[![goreportcard for breathbath/gotainer](https://goreportcard.com/badge/github.com/breathbath/gotainer?v=1)](https://goreportcard.com/report/breathbath/gotainer)
[![codecov for breathbath/gotainer](https://codecov.io/gh/breathbath/gotainer/branch/master/graph/badge.svg)](https://codecov.io/gh/breathbath/gotainer)
[![Sourcegraph for breathbath/gotainer](https://sourcegraph.com/github.com/breathbath/gotainer/-/badge.svg)](https://sourcegraph.com/github.com/breathbath/gotainer?badge)

This library helps to manage dependencies in your project by providing a centralised logic for initialising services.

You can define any go lang type as a service - a struct, a function (closure or lambda), a scalar value or a dynamic result of your
functions. You do it in a simple manner: you create a container instance and add your services to it under an unique alias.
Further on you can fetch them by this key in any part of your application.

# How to install

        go get github.com/breathbath/gotainer/container

If you use [Go dep tool](https://github.com/golang/dep):

        //1. Add it to your dependencies file
        dep ensure -add github.com/breathbath/gotainer/container

        //2. Use container somewhere in the code (e.g. declare some dependencies)

        //3. Fix the result
        dep ensure

# Quick start

## Declaring services

        //first we create a container
        container := container.NewRuntimeContainer()

        //then we declare a callback that will return MyService instance identified by "my_service"
        runtimeContainer.AddConstructor("my_service", func(c container.Container) (interface{}, error){
            return MyService{}, nil
        })

If you already have a constructor function, you can add it to the container as well:

        func NewMyService() MyService {
            return MyService{}
        }

        runtimeContainer.AddNewMethod("my_service", NewMyService)

If you want to create a config file for your dependencies, use:

        configTree := Tree{
        		Node{
        			NewFunc: mocks.NewConfig,
        			Id:      "config",
        		},
        		Node{
        			Id: "connection_string",
        			Constr: func(c Container) (interface{}, error) {
        				config := c.Get("config", true).(mocks.Config)
        				return config.GetValue("fakeDbConnectionString"), nil
        			},
        		},
        		Node{
        			Id:           "db",
        			NewFunc:      mocks.NewFakeDb,
        			ServiceNames: Services{"connection_string"},
        		},
        		Node{
        			Id: "authors_storage_statistics_provider",
        			Ev: Event{
        				Name:    "add_stats_provider",
        				Service: "authors_storage",
        			},
        		},
        		Node{Id: "statistics_gateway", NewFunc: mocks.NewStatisticsGateway},
        		Node{
        			Ob: Observer{
        				Event: "add_stats_provider",
        				Name:  "statistics_gateway",
        				Callback: func(sg *mocks.StatisticsGateway, sp mocks.StatisticsProvider) {
        					sg.AddStatisticsProvider(sp)
        				},
        			},
        		},
        	}

        ...
        builder := RuntimeContainerBuilder{}
        //at this point you have a fully working container with dependencies from the config tree
        container := builder.BuildContainerFromConfig(configTree)

Please note that you can declare your services in any order, e.g. you can declare a service A which is dependent
on service B before or after it.

## Fetching services
Assuming that you already created a container and declared all needed services, you can start fetching them:

        var myService MyService
        container.Scan("my_service", &myService)
        //at this point myService will contain the initialised instance of MyService, which was either created by
        //the provided callback constructor or by your custom New function
        myService.SomeMethod()

# Use cases

The library covers the following use cases:

## Reusable services with dependencies

Imagine we have 2 services dependant on another one, which should be reused in both cases. The common service
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

4. You might have the container declaration for your MonitoringGateway in one core library and different implementations of
MonitoringProvider in other packages, so you are able to plug them in individually in every application with no need to change the
core code.

## Parameters

Parameters are simple scalar values, that are already defined and can be used as dependencies. A typical example is an application config with
e.g. db connection details, proxy url or folder paths. Obviously you can have a service that requires parameters in a construction
function. Gotainer provides a useful helper function, that allows to add parameters as a `map[string] interface{}`.
Consider following example:

        //e.g. map[string] string {"password": "123456", "proxy": "127.0.0.1:8888"}
        type Config map[string] string
        ...
        //some custom function that converts json to Config type
        config := fetchFromFile("config.json")
        ...
        container := createContainer()

        //this adds each map key/value pair as a single parameter, so each of it can be addressed by it's name
        RegisterParameters(container, config)

        //proxyConnectionStr is now "127.0.0.1:8888"
        proxyConnectionStr := container.Get("proxy").(string)

        //or you can use it as a dependency
        c.AddConstructor("pass_checker", func(c Container) (interface{}, error) {
            var password string
            c.Scan("password", &password)
            //some custom class that checks if password provided by user is correct
            return PasswordChecker{correctPassword: password), nil
        })

        //or like this
        c.AddNewMethod("pass_checker", NewPassChecker, "password")

You can also declare parameters in a dependency config:

        //statically
        ...
        Node {
            Parameters: map[string]interface{}{
                "param1": "value1",
                "param2": 123,
            },
        },
       ...

       //or dynamically via interface
        type ConfigProvider struct{}

        //implements container.ParametersProvider
        func (cp ConfigProvider) GetItems() map[string]interface{} {
            return map[string]interface{}{
                "EnableLogging": true,
            }
        }

        ...
        Node {
            ParamProvider: ConfigProvider{},
        },
        ...

# Good practices

## Creating a dependency container

1. Declare a function that will be responsible for the container creation, e.g.

        package app_container

        func NewAppContainer() RuntimeContainer {
            container := container.NewRuntimeContainer()

            //services declarations here...

            return container
        }

If your application has other libraries that use the container, you can merge all dependency declarations into one.
If your application is big, you can declare small containers for your packages merge them in your main container method.

        package app_container

        import other_library_container "github.com/myname/other_library/container"

        func NewAppContainer() RuntimeContainer {
            container := container.NewRuntimeContainer()

            //services declarations here...

            otherLibraryContainer := other_library_container.NewAppContainer()
            container.Merge(otherLibraryContainer)

            return *container
        }

You can also merge dependencies configs into one result container like this:

        treeFromModule1 := NewModuleOneConfigTree() //imagine this method returns container.Tree{}
        treeFromModule2 := NewModuleTwoConfigTree()
        //at this point your container will have dependencies from both treeFromModule1 and treeFromModule2
        container := RuntimeContainerBuilder{}.BuildContainerFromConfig(treeFromModule1, treeFromModule2)

Don't put container init logic into your main.go file as it might grow very big and will not be reusable.

The best way to avoid this is to return a container from your "NewAppContainer" method rather than a pointer to it.
This will make sure that your container won't be modified at runtime in your business code.

2. Add services declarations in the container init method:

        package app_container

        func NewAppContainer() RuntimeContainer {
            container := container.NewRuntimeContainer()

           runtimeContainer.AddNewMethod("service_1", NewService1)
           runtimeContainer.AddNewMethod("service_2", NewService2, "service_1")
           runtimeContainer.AddNewMethod("service_3", NewService2, "service_1", "service_2")

            return container
        }

3. If you have services with optional dependencies, declare them via callbacks:

        //...
        runtimeContainer.AddConstructor("service_a", func(c container.Container) (interface{}, error) {
            var logger Logger
            c.Scan("logger", &logger)

            myService := MyService{}
            myService.SetLogger(logger)

            return myService
        })
        //...

4. Don't declare container as a dependency for a service.

        type MyType struct {
            container Container
        }


Generally it's a bad practise for the following reasons:

- Unit testing of such services will be cumbersome as you would need to mock an undefined amount of dependencies,
that your code might require from the container

- Dependencies for your service will be hidden inside, so its public interface will be less obvious for understanding

- You couple your code with the container library, which may produce an overhead in it's usage in other applications or projects

- You run the risk of producing circular dependencies (e.g., your service asks the container for a dependency which requires your service)

## Fetching services

1. Fetch your services only in the main.go method.

       package main

       func main() {
            container := NewAppContainer()

            var interactor SomeInteractor
            container.Scan("interactor", &interactor)

            interactor.DoSomething()
       }

2. Don't pass container as a dependency to your business logic, only your controllers should communicate to it.

3. Use "scan" methods to get typed services. Don't forget to use pointer types in the destination argument (otherwise there will be a panic error).
You can use interface return types and type assertions as well like this:

       package main

       func main() {
            container := NewAppContainer()

            interactor := container.Get("interactor").(Interactor)

            interactor.DoSomething()
       }

## Testing

Working with "RuntimeContainer" means that possible errors in a service declaration won't appear until you fetch it from the container.
To make sure, that your declared container has valid service definitions, you should run the "Check" method. You
do it in an integration test as:

        func TestContainer(t *testing.T) {
            container := NewAppContainer()
            container.Check()
        }

## Dependencies cache
Dependencies cache is a in-memory storage for dependencies allowing to retrieve a service in an initialized state. 
This gives a great opportunity to share services between different consumers to spare time for initialisation. 
A typical example for that is a db connection, which can be reused by different db dependant services, without 
explicitly opening connection each time you need to query one or another table/objects collection.

RuntimeContainer allows you to get a cached or non-cached version of a service through different methods:

	Scan(id string, dest interface{}) //scans a cached version of a service into a destination
	ScanNonCached(id string, dest interface{}) //scans a non-cached version
	Get(id string, isCached bool) interface{} //receive service in the output where isCached flag affects the cache switch
	GetSecure(id string, isCached bool) (interface{}, error) //the same with error driven version of Get
	
In any case, if you require a non cached version of a service, it will be initialised from beginning and will be cached
replacing the old version.

An alternative non-container approach is to use the go's init() function where you can put your initialisation logic, 
which will be executed also only once. But in this case you should implement your own logic to get an uncached version of a service 
which is sometimes cumbersome and also quite repetative. With the RuntimeContainer you get this opportunity out of the box.

## Lazy initialisation

If you use the RuntimeContainer, you have an advantage of a lazy services initialisation which is provided out of the box.
This means that any service you declare in the container is not initialised immediately but when it is explicitly required
directly or by some other service. This means practically that your application won't open a db connection, when your
current execution is not requiring a db data (meaning a db dependent service).

If you use an init() function for a service initialisation, this will be always executed, once it's package is imported
somewhere. To have a lazy initialisation you should create your own builder methods which will be a boilerplate code,
which we want to avoid with a dependency container.


## Garbage collection

Sometimes your code might use resources which should be released on the application exit. One typical example is a db connection
which is expected to be closed, when your db client is ready with all operations. The recommended and obvious way is to use the defer operator e.g.

    defer dbConn.close()

Imagine that the dbConn is reused by many different services you define in your application:

    dbConn := db.Connect()
    usersProvider := users.NewProvider(dbConn)
    articlesProvider :=  articles.NewProvider(dbConn)

    defer dbConn.close()

That is perfectly ok, but if you decided to use a dependency container to get rid of a boilderplate code, you need to delegate the releasing of
resources (garbage collection) to the container.

    cont := NewAppContainer()
    articlesProvider := cont.get("articles_provider", true).(ArticlesProvider)
    //dbConn is injected by the container into articlesProvider so you don't have access to it here

One might think to release a shared resource in on of the services using it. But what if another service later calls the dbConn and
discovers it's closed? This will lead to very hard trackable errors.

The Gotainer has a garbage collection functionality to solve this problem. You can register all your garbage collection functions and call them just by doing:

    defer container.CollectGarbage()

Let's look at some examples:

       ...

       //my db service holds reference to the db connection
       type MyDbService struct{
            dbConn: DbConnection
       }

       //imagine some library method returning a DbConnection instance
       func NewMyDbService(dbConn: DbConnection) MyDbService {
          return MyService{dbConn: dbConn}
       }

       //this is the actual garbage collection
       func (ms MyService) Destroy() error {
           return ms.dbConn.Disconnect()
       }

       ...
       func main() {
            container := NewAppContainer() //here we already have created a DbConnection resource
            container.AddNewMethod("my_db_service", NewMyDbService, "dbConnString") //adding our services with the dbConnString in it

            //the actual garbage collection callback
            dbGarbageCollector := func(service interface{}) error {
                return service.(MyDbService).Destroy()
            }
            container.AddGarbageCollectFunc("my_db_service", dbGarbageCollector)

            //this will happen on exit from main (and therefore from the application rather than in some service)
            defer container.CollectGarbage()

            //or you can also notify about disconnection problems
            defer func() {
                err := container.CollectGarbage()
                if err != nil {
                    fmt.Println(err)
                }
            }
       }

If you prefer to use container config:

        	return Tree{
        	    ...
        		Node{
        			Id:           "my_db_service",
        			NewFunc:      NewMyDbService,
        			ServiceNames: Services{"dbConnString"},
        			GarbageFunc: func(service interface{}) error {
        				myDbService := service.(MyDbService)
        				return myDbService.Destroy()
        			},
        		},
        		...
        	}

In this case just don't forget to call `defer container.CollectGarbage()` in your main function.

Garbage collection functions are called in the order of declaration. You should avoid calling release function of already released resources.
In this case the shared services should register just one garbage collection call, rather than the services using it trying to release them in their own
release functions.

Let's consider the following case:

        		Node{
        			Id:           "my_db_service",
        			NewFunc:      NewMyDbService,
        			ServiceNames: Services{"dbConnString"},
        			GarbageFunc: func(service interface{}) error {
        				myDbService := service.(MyDbService)
        				return myDbService.Destroy()
        			},
        		},
                Node{
                    Id:           "users_provider",
                    NewFunc:      NewUsersProvider,
                    ServiceNames: Services{"my_db_service"},
                    GarbageFunc: func(service interface{}) error {
                        myUsersProvider := service.(UsersProvider)
                        return myUsersProvider.Destroy()
                    },
                },

                type UsersProvider struct {
                    db MyDbService
                    externalConn ExternalConn
                }

                func NewUsersProvider(db MyDbService) UsersProvider {
                    return UsersProvider {db: db, externalConn: externalLibrary.OpenConn()}
                }

                func (up UsersProvider) Destroy() error {
                    up.externalConn.close() //here we release resources owned by the UsersProvider only which is ok
                    up.db.Destroy() //don't do this as this will be called earlier by the container

                    return nil
                }

The general rule is that shared services are responsible for garbage collection calls, rather than services using them.

## Cycle detection
Dependency cycle is a classic problem of graph cycles in [computer science](https://en.wikipedia.org/wiki/Cycle_(graph_theory)).
A cycle is a situation where one dependency requires itself as a constructor argument or appears in the requirement list
of any other dependency which is needed to construct the current one. 
Here are some examples: 

A self reference cycle (db -> db):

    ...
    //a db service is fetching itself upon construction
	cont.AddConstructor("db", func(c Container) (interface{}, error) {
		return c.GetSecure("db")
	})
	...
	
A dependencies circle (rolesProvider -> userProvider -> rolesProvider):
	
		cont.AddConstructor("roleProvider", func(c Container) (interface{}, error) {
    		return c.GetSecure("userProvider", true)
    	})
    
    	cont.AddConstructor("userProvider", func(c Container) (interface{}, error) {
    		return c.GetSecure("roleProvider", true)
    	})
    	
The same with in the config declaration

        Tree{
            Node{
                Id:           "userProvider",
                NewFunc:      mocks.NewUserProvider,
                ServiceNames: Services{"roleProvider"},
            },
            Node{
                Id:           "roleProvider",
                NewFunc:      mocks.NewRoleProvider,
                ServiceNames: Services{"userProvider"},
            }
        }

If cycle appears, the dependency cannot be created which leads to a runtime error. In this case if you use `Scan/Get/ScanNonCached` methods, 
application will panic. If you're using `ScanSecure/GetSecure` methods, you will get a cycle detection error in the method output result.

Detection of a dependency cycle means that the container is unusable and a developer must change the code to fix this error. 
Unfortunately Runtime container cannot detect cycles at compile time. This is the price you should pay for a lazy 
dependencies initialisation. 
To detect cycles in a container, you should trigger the `Check` function e.g. in a test. 
We recommend to create a simple test as mentioned [here](https://github.com/breathbath/gotainer#testing) and setup a CI env
to detect cycles before the faulty code goes to production.