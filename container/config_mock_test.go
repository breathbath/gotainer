package container

import "github.com/breathbath/gotainer/container/mocks"

func getMockedConfigTree() Tree {
	return Tree{
		Node{
			ID: "connection_string",
			Constr: func(c Container) (interface{}, error) {
				config := c.Get("config", true).(mocks.Config)
				return config.GetValue("fakeDbConnectionString"), nil
			},
		},
		Node{
			ID:           "book_storage",
			NewFunc:      mocks.NewBookStorage,
			ServiceNames: Services{"db"},
		},
		Node{
			NewFunc: mocks.NewConfig,
			ID:      "config",
		},
		Node{
			ID: "book_creator",
			Constr: func(c Container) (interface{}, error) {
				return mocks.BookCreator{}, nil
			},
		},
		Node{
			ID:           "book_finder_declared_statically",
			NewFunc:      mocks.NewBookFinder,
			ServiceNames: Services{"book_storage", "book_creator"},
		},
		Node{
			ID: "book_storage_statistics_provider",
			Ev: Event{
				Name:    "add_stats_provider",
				Service: "book_storage",
			},
		},
		Node{
			ID: "book_finder",
			Constr: func(c Container) (interface{}, error) {
				var bc mocks.BookCreator
				c.Scan("book_creator", &bc)

				var bs mocks.BookStorage
				c.Scan("book_storage", &bs)

				return mocks.NewBookFinder(bs, bc), nil
			},
		},
		Node{
			ID:           "db",
			NewFunc:      mocks.NewFakeDb,
			ServiceNames: Services{"connection_string"},
			GarbageFunc: func(service interface{}) error {
				fakeDb := service.(*mocks.FakeDb)
				return fakeDb.Destroy()
			},
		},
		Node{
			ID:           "authors_storage",
			NewFunc:      mocks.NewAuthorsStorage,
			ServiceNames: Services{"db"},
		},
		Node{
			ID: "authors_storage_statistics_provider",
			Ev: Event{
				Name:    "add_stats_provider",
				Service: "authors_storage",
			},
		},
		Node{ID: "statistics_gateway", NewFunc: mocks.NewStatisticsGateway},
		Node{
			Ob: Observer{
				Event: "add_stats_provider",
				Name:  "statistics_gateway",
				Callback: func(sg *mocks.StatisticsGateway, sp mocks.StatisticsProvider) {
					sg.AddStatisticsProvider(sp)
				},
			},
		},
		Node{
			ID: "static_files_url",
			Constr: func(c Container) (interface{}, error) {
				var config mocks.Config
				c.Scan("config", &config)

				return config.GetValue("staticFilesUrl"), nil
			},
		},
		Node{
			ID:           "book_link_provider",
			NewFunc:      mocks.NewBookLinkProvider,
			ServiceNames: Services{"static_files_url", "book_finder_declared_statically"},
		},
		Node{ID: "web_fetcher", NewFunc: mocks.NewWebFetcher},
		Node{ID: "in_memory_cache", NewFunc: mocks.NewInMemoryCache},
		Node{
			Parameters: map[string]interface{}{
				"param1": "value1",
				"param2": 123,
			},
			ParamProvider: mocks.ConfigProvider{},
		},
		Node{
			ID:           "logger",
			NewFunc:      mocks.BuildLogger,
			ServiceNames: Services{"EnableLogging"},
		},
	}
}

func buildContainerFromMockedConfig() Container {
	configTree := getMockedConfigTree()
	return RuntimeContainerBuilder{}.BuildContainerFromConfig(configTree)
}
