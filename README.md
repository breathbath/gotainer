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

Please note that you can declare services in any order, e.g., first A which requires B and then B itself. B can be
also declared in another config file and then merged with the definitions containing A.

## Fetching services
After you created a container and declared all needed services, we can start fetching them:

        var myService MyService
        container.Scan("my_service", &myService)
        //at this point myService will contain the initialised instance of MyService, which was either created by
        //the provided callback constructor or by your custom New function

        //or using interface assertions
        myService := container.Get("my_service", true).(MyService)
        //from this point the service is fully functional
        myService.SomeMethod()

# Use cases

## Shared states

Imagine that we have 2 services dependant on the same one, which should be reused in both cases. The common service
should be initialised only once as it has an internal state (e.g., a db connection).

        type ServiceA struct {dbConn *DbConn}

        type ServiceB struct {
            serviceA ServiceA
        }

        type ServiceC struct {
           serviceA ServiceA
        }

        container.AddNewMethod("service_a", NewServiceA)
        container.AddNewMethod("service_b", NewServiceB, "service_a")
        container.AddNewMethod("service_c", NewServiceC, "service_a")

        //when fetching "service_b" or "service_c" container will make sure that they share "service_a" in an
        //initialised shared state e.g. with an open and reusable db connection
        var serviceC ServiceC
        container.Scan("service_c", &serviceC)
        serviceB := container.Get("service_b", true).(ServiceB)

## Services with a complex initialisation

        type ServiceX struct {...}
        func (sx ServiceX) AddService(sc ServiceC){}
        func (sx ServiceX) EnableLogging(){}
        func (sx ServiceX) RegisterInMonitoringList(monitoringList []MonitoringItem){}
        container.AddConstructor("service_x",  func(c container.Container) (interface{}, error){
            serviceX := NewServiceX()

            var serviceC ServiceC
            c.Scan("service_c", &serviceC)
            serviceX.AddService(serviceC)

            serviceX.EnableLogging()

            var monitoringList []MonitoringItem
            c.Scan("monitoring_list", &monitoringList)
            serviceX.RegisterInMonitoringList(monitoringList)

            return serviceX, nil
        })

We're doing a complex ServiceX initialisation. This code will be executed once.
All services using "service_x" will have the correctly initialised version of it.

## Cached and reusable results of a func execution
        //in this case SumItems can be an expensive operation that should be executed once
        func SumItems(items []SomeItems) int64 {
            var result int64 = 0
            for _, item := range items {
                result += item.Sum()
            }
            return result
        }

        container.AddConstructor("items_sum", func(c container.Container) (interface{}, error) {
            itemsProvider := c.Get("items_provider", true).(ItemsProvider)
            return SumItems(itemsProvider.GetItems()), nil
        })

        //now we can call this 100 times in different places but SumItems will be executed only once:
        container.Get("items_sum").(int64)

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

## Anonymous constructors for simply initialised services

        container.AddConstructor("service_a", func(c container.Container) (interface{}, error) {
               return ServiceA{}
        })

## Dependency events

In many cases your service wants get dependencies of a certain type every time when they are added to the container but it should
stay unmodified (see [Open-Closed-Principle](https://en.wikipedia.org/wiki/Open%E2%80%93closed_principle)).
Consider following example:

        type MonitoringProvider interface{
            GetMonitoringEvent() (eventName string, count int64)
        }

        //first implementation of MonitoringProvider
        type TotalMonitoringProvider struct{...}
        func(tmp TotalMonitoringProvider) GetMonitoringEvent{
            return "total_count", 100
        }

        //second implementation of MonitoringProvider
        type ErrorCountProvider struct{...}
        func(tmp ErrorCountProvider) GetMonitoringEvent{
            return "errors_count", 10
        }

        //wants to have all implementations of MonitoringProvider
        type MonitoringGateway struct{monitorigProviders []MonitoringProvider}
        func (mg MonitoringGateway) AddMonitoringProvider (mp MonitoringProvider)...


In this case we want to add all existing and future implementations of MonitoringProvider to the MonitoringGateway.
An obvious solution could be:

        func BuildMonitoringGateway(tmp TotalMonitoringProvider, ecp ErrorCountProvider) MonitoringGateway {
            mg := MonitoringGateway{}
            mg.AddMonitoringProvider(tmp)
            mg.AddMonitoringProvider(ecp)
        }

This approach has following problems:

1. With every new implementation of MonitoringProvider you should modify the BuildMonitoringGateway, so this code might break.

2. The amount of arguments of BuildMonitoringGateway will grow, so this function becomes unreadable.

3. You should create every new instance of MonitoringProvider somewhere which will probably lead to code duplication, if those
require other services, the amount of boilerplate code will explode

With the Gotainer you can solve this problem as following:

        //we declare monitoring gateway the NewMonitoringGateway is free from any dependencies
        container.AddNewMethod("monitoring_gateway", NewMonitoringGateway)
        // we say here that "monitoring_gateway" is interested in the event "monitoring_provided_added", and every time it happens
        //we should execute function func(mg MonitoringGateway, mp MonitoringProvider) which adds every new implementation of MonitoringProvider to the
        //monitoring gateway
        container.AddDependencyObserver("monitoring_provided_added", "monitoring_gateway", func(mg MonitoringGateway, mp MonitoringProvider){
            sg.AddMonitoringProvider(sp)
        })

        container.AddNewMethod("total_monitoring_provider", NewTotalMonitoringProvider, "service_a", "service_b")
        container.RegisterDependencyEvent("monitoring_provided_added", "total_monitoring_provider")

        container.AddNewMethod("error_count_provider", NewErrorCountProvider)
        container.RegisterDependencyEvent("monitoring_provided_added", "error_count_provider")

        //we can add further other possible provides without changing code of MonitoringGateway

This code has following advantages:

1. MonitoringGateway is decoupled from new implementations of MonitoringProvider.

2. No complex initialisation function for MonitoringGateway is needed.

3. Concrete implementations of MonitoringProvider are created once without any repetition as this logic is already encapsulated in the Gotainer.

4. You might have the container declaration for your MonitoringGateway in one core library and different implementations of
MonitoringProvider in other packages, so you are able to plug them in individually in every application with no need to change the
core code.

## Shared application parameters

Parameters are simple scalar values, that are defined in config files and can be used as dependencies. A typical example is an application config with
e.g. db connection details, proxy url or folder paths. Obviously you can have a service that requires parameters in a construction
function. Gotainer provides a useful helper function, that allows to add parameters as a `map[string] interface{}`.
Consider following example:

        //e.g. map[string] string {"password": "123456", "proxy": "127.0.0.1:8888", "url": "www.domain.com", "is_log_enabled": true, "max_failures_count": 2}
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

        //now all parameters can be used as usual shared services
        container.AddNewMethod("api_caller", NewApiCaller, "url", "is_log_enabled", "max_failures_count")

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

## Creating the dependency container

1. Declare a function that will be responsible for the container creation, e.g.

        package app_container

        func NewAppContainer() RuntimeContainer {
            container := container.NewRuntimeContainer()

            //services declarations will go here...

            return container
        }

If your application has other libraries that use the container, you can merge all dependency declarations into one.
If your application is big, you can declare small containers for packages, merge them in your main container method.

        package app_container

        import other_library_container "github.com/myname/other_library/container"

        func NewAppContainer() RuntimeContainer {
            container := container.NewRuntimeContainer()

            //services declarations will go here...

            otherLibraryContainer := other_library_container.NewAppContainer()
            container.Merge(otherLibraryContainer)

            return *container
        }

You can also merge configs into one container like this:

        treeFromModule1 := NewModuleOneConfigTree() //imagine this method returns container.Tree{}
        treeFromModule2 := NewModuleTwoConfigTree()
        //at this point your container will have dependencies from both treeFromModule1 and treeFromModule2
        container := RuntimeContainerBuilder{}.BuildContainerFromConfig(treeFromModule1, treeFromModule2)

Don't put container init logic into your main.go file as it might become very big and unreadable.

The best way to avoid this, is to return the container from your "NewAppContainer" method rather than a pointer to it.
This will make sure that your container won't be modified at runtime in your business code.

2. Add services declarations in the container's build method:

        package app_container

        func NewAppContainer() RuntimeContainer {
           container := container.NewRuntimeContainer()

           runtimeContainer.AddNewMethod("service_1", NewService1)
           runtimeContainer.AddNewMethod("service_2", NewService2, "service_1")
           runtimeContainer.AddNewMethod("service_3", NewService2, "service_1", "service_2")

           return container
        }

3. If you have services with optional dependencies, declare them as anonymous constructors:

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

2. Don't pass container as a dependency to your business logic, only your entry points should communicate to it.

3. Use "scan" methods to get typed services. Don't forget to use pointer types in the destination argument (otherwise there will be a panic error).
You can use interface type assertions as well:

       package main

       func main() {
            container := NewAppContainer()

            interactor := container.Get("interactor").(Interactor)

            interactor.DoSomething()
       }

## Testing

Working with "RuntimeContainer" means that possible errors in a service declaration won't appear until you fetch it from the container.
To make sure, that your declared container has valid service definitions, you should run the "Check" method:

        func TestContainer(t *testing.T) {
            container := NewAppContainer()
            container.Check()
        }

It's not recommended to use this method in production env, as it might lead to performance issues. The "Check" method
requires creation of all services from the container, which of course can be avoided for many use cases.
Imagine that you trigger a db migration function which requires to initialise only the database service. Calling "Check" here
will initialise a lot of not needed services. 
On the other hand using testing allows to validate all container dependencies before pushing code to production.

## Dependencies cache
Dependencies cache is a in-memory storage allowing to retrieve a service in an initialized shareable state.
This gives a great opportunity to share services among different consumers to spare time for initialisation.
A typical example for that is a db connection, which can be used by different db related services, without
explicitly opening it each time you need to query some table.

RuntimeContainer allows you to get a cached or a non-cached version of a service by different methods:

	Scan(id string, dest interface{}) //scans a cached version of a service into a destination
	ScanNonCached(id string, dest interface{}) //scans a non-cached version
	Get(id string, isCached bool) interface{} //receive service in the output where isCached flag affects the cache switch
	GetSecure(id string, isCached bool) (interface{}, error) //the same with error driven version of Get
	
In any case, if you require a non cached version of a service, it will be initialised from beginning and will be cached
replacing the old version.

An alternative non-container approach is to use the go's init() function where you can put your initialisation logic, 
which will be executed also only once. But in this case you should implement your own logic to get an uncached version of a service 
which is sometimes cumbersome and also quite repetative. The RuntimeContainer provides this functionality out of the box.

## Lazy initialisation

RuntimeContainer gives advantage of lazy services initialisation. This means that any declared service is created only when it is explicitly requested from the
container or another service needs it for it's creation. This allows to spare resources by creating only the subset of services needed for the current execution.

Using init() function for initialisation would mean just the opposite - once the module is loaded, it is executed even if it
is not used in the current call. Creating own lazy initialisation is a cumbersome process where code repetition is hard to avoid.

The RuntimeContainer provides this functionality out of the box.

## Garbage collection

Sometimes your code might use resources which should be released on the application exit. One typical example is a db connection
which is expected to be closed, when your db client has finished all operations. The recommended and obvious way is to use the defer operator e.g.

    defer dbConn.close()

Imagine that the dbConn is reused by many different services in your application:

    dbConn := db.Connect()
    usersProvider := users.NewProvider(dbConn)
    articlesProvider :=  articles.NewProvider(dbConn)

    defer dbConn.close()

But what if you return from the current method but some service is still using the db connection? Of course it would be nice to release resources
only if the main() method finishes it's work.

You can delegate the releasing of resources (garbage collection) to the container.

    //we initialise the container in the main method
    cont := NewAppContainer()
    //will always be executed on after the main method, now all services which registered a GarbageCollect function will get released
    defer cont.CollectGarbage()

Let's look at examples:

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
In this case the shared resources should register just one garbage collection call, rather than the services using it trying to release them in their own
release functions.

Let's consider the following case:

        	    Node{
                    Id: "my_db_service",
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
Dependency cycle is a classic case of graph cycles in [computer science](https://en.wikipedia.org/wiki/Cycle_(graph_theory)).
A cycle is a situation where one dependency requires itself as a constructor argument or appears in the requirement list
of any other dependency which is needed to construct the current one. 
Here are some examples: 

Self reference (db -> db):

    ...
    //a db service is fetching itself upon construction
	cont.AddConstructor("db", func(c Container) (interface{}, error) {
		return c.GetSecure("db")
	})
	...
	
Dependencies circle (rolesProvider -> userProvider -> rolesProvider):
	
		cont.AddConstructor("roleProvider", func(c Container) (interface{}, error) {
    		return c.GetSecure("userProvider", true)
    	})
    
    	cont.AddConstructor("userProvider", func(c Container) (interface{}, error) {
    		return c.GetSecure("roleProvider", true)
    	})
    	
The same in the config declaration:

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

If a cycle appears, the dependency cannot be created which leads to a runtime error. In this case if you use `Scan/Get/ScanNonCached` methods, 
application will panic. If you're using `ScanSecure/GetSecure` methods, you will get a cycle detection error in the method output result.

If the dependency cycle is detected, the container becomes unusable and the developer must change the container definition to fix the error.
Unfortunately in the RuntimeContainer we cannot detect cycles at compile time. This is the price you should pay for a lazy
dependencies initialisation. 

To detect possible cycles in a container, you should trigger the container's `Check` function in a go test. 
We recommend to create a simple test as mentioned [here](https://github.com/breathbath/gotainer#testing) and setup a CI env
to detect cycles before the faulty code goes to production.

## Duplicates detection
It's very common that service declarations are copy pasted. But sometimes you might forget to change the pasted service id. 
The default behaviour of the RuntimeContainer is that the latest values duplicate the previous ones.

For example you declare something like this:

    	return Tree{
    		Node{
    			ID:           "proxyProvider",
    			NewFunc:      proxy.ProxyProvider,
    		},

    		... //many other services

    	}

then you copy paste the code to the bottom but forget to change the ID:

    	return Tree{
    		Node{
    			ID:           "proxyProvider",
    			NewFunc:      proxy.NewProxyProvider,
    		},

    		... //many other services

            Node{
                ID:           "proxyProvider",
                NewFunc:      io.NewFileManager,
            },
    	}

After you fetch the "proxyProvider" from the container, you get a nasty error that the "proxyProvider" service is not compatible with the ProxyProvider type.
You check the config again, find the "proxyProvider" on the top line and see no reason for the problem as `proxy.NewProxyProvider` returns ProxyProvider.
Everything looks fine but the error is still there. Sometimes finding the reason for such errors can take time especially for long or split container
configs.

To avoid such kind of problems, a duplicates detector was added to the RuntimeContainer. In the case above container creation will fail with a meaningful error
text.

The duplicates detector is triggered every time you use `AddConstructor` or `AddNewMethod` of the `RuntimeContainer` or when building container from the config.

If you want to declare services without the duplicates detector (e.g., for overriding existing services), you might use the `SetNewMethod` and `SetConstructor`
functions. More about this see [here](https://github.com/breathbath/gotainer#overriding-existing-services)

## Overriding existing services
Overriding existing container definitions is helpful in testing environments, where you want to replace resource critical or "production only" services with their mocked implementations.

A typical example might be a service which calls an external API e.g., a payment gateway. Of course in testing environment we want to avoid
creating real payments for our test scenarios. The problem is easily solved by replacing the payment service with a dummy code, which mimics real API responses.

You can find an example for this use case in [this test](https://github.com/breathbath/gotainer/blob/master/container/container_test.go#L177)
as well as [here](https://github.com/breathbath/gotainer/blob/master/container/mocks/paymentGateway.go).

It's hard to imagine a valid use case for using services overriding in a prod environment. The complexity of initialisation logic
should be encapsulated in the container rather than in the code outside.
You should remember that the container is not a part of any business logic. It's just a helper which reduces the boilerplate code
and arranges initialization logic in one place.

Doing something like:


        cont := BuildMyAppContainer()
        if os.Getenv("IS_PROXY_DISABLED") == "1" {
            cont.SetNewMethod("proxy", NewNullProxy)
        }


is a bad idea, because initialisation code is outside of the container and the unnecessary logic is exposed to the its caller.
Accordingly our "business" code knows too much about different ways to create and replace the `proxy`.

Imagine that the `BuildMyAppContainer` func is reused in many places, but only here we get different implementation for the `proxy` service.
You probably would prefer to have this logic everywhere we call `BuildMyAppContainer`.

We improve the code by modifying the `BuildMyAppContainer` func like this:


        func BuildMyAppContainer() *container.RuntimeContainer {
            cont := container.NewRuntimeContainer()
            cont.AddConstructor("proxy", func(c container.Container) (i interface{}, e error) {
                if os.Getenv("IS_PROXY_DISABLED") == "1" {
                    return proxy.NewNullProxy()
                }

                return proxy.NewRealProxy()
            })
        }

Now we're hiding the proxy initialisation details in the place where we create the container and making sure that all
calls of `BuildMyAppContainer` will provide the same logic which is a preferable approach.