package container

import "github.com/breathbath/gotainer/container/mocks"

func getMockedConfigTree() Tree {
	return Tree{
		Node{
			Id: "connection_string",
			Constr: func(c Container) (interface{}, error) {
				config := c.Get("config", true).(mocks.Config)
				return config.GetValue("fakeDbConnectionString"), nil
			},
		},
		Node{
			Id:           "book_storage",
			NewFunc:      mocks.NewBookStorage,
			ServiceNames: Services{"db"},
		},
		Node{
			NewFunc: mocks.NewConfig,
			Id:      "config",
		},
		Node{
			Id: "book_creator",
			Constr: func(c Container) (interface{}, error) {
				return mocks.BookCreator{}, nil
			},
		},
		Node{
			Id:           "book_finder_declared_statically",
			NewFunc:      mocks.NewBookFinder,
			ServiceNames: Services{"book_storage", "book_creator"},
		},
		Node{
			Id: "book_storage_statistics_provider",
			Ev: Event{
				Name:    "add_stats_provider",
				Service: "book_storage",
			},
		},
		Node{
			Id: "book_finder",
			Constr: func(c Container) (interface{}, error) {
				var bc mocks.BookCreator
				c.Scan("book_creator", &bc)

				var bs mocks.BookStorage
				c.Scan("book_storage", &bs)

				return mocks.NewBookFinder(bs, bc), nil
			},
		},
		Node{
			Id:           "db",
			NewFunc:      mocks.NewFakeDb,
			ServiceNames: Services{"connection_string"},
		},
		Node{
			Id:           "authors_storage",
			NewFunc:      mocks.NewAuthorsStorage,
			ServiceNames: Services{"db"},
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
		Node{
			Id: "static_files_url",
			Constr: func(c Container) (interface{}, error) {
				var config mocks.Config
				c.Scan("config", &config)

				return config.GetValue("staticFilesUrl"), nil
			},
		},
		Node{
			Id:           "book_link_provider",
			NewFunc:      mocks.NewBookLinkProvider,
			ServiceNames: Services{"static_files_url", "book_finder_declared_statically"},
		},
		Node{Id: "web_fetcher", NewFunc: mocks.NewWebFetcher},
		Node{Id: "in_memory_cache", NewFunc: mocks.NewInMemoryCache},
		Node {
			Parameters: map[string]interface{}{
				"param1": "value1",
				"param2": 123,
			},
			ParamProvider: mocks.ConfigProvider{},
		},
		Node {
			Id: "logger",
			NewFunc:      mocks.BuildLogger,
			ServiceNames: Services{"EnableLogging"},
		},
	}
}
