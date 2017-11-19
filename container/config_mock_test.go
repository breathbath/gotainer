package container

import "github.com/breathbath/gotainer/container/mocks"

func getMockedConfigTree() Tree {
	return Tree{
		"config": Node{
			NewFunc: mocks.NewConfig,
		},
		"connection_string": Node{
			Constr: func(c Container) (interface{}, error) {
				return c.Get("Tree", true).(string), nil
			},
		},
		"db": Node{NewFunc: mocks.NewFakeDb, ServiceNames: Services{"connection_string"}},
		"book_finder_declared_statically": Node{
			NewFunc:      mocks.NewBookFinder,
			ServiceNames: Services{"book_storage", "book_creator"},
		},
		"book_finder": Node{
			Constr: func(c Container) (interface{}, error) {
				var bc mocks.BookCreator
				c.Scan("book_creator", &bc)

				var bs mocks.BookStorage
				c.Scan("book_storage", &bs)

				return mocks.NewBookFinder(bs, bc), nil
			},
		},
		"book_storage": Node{
			NewFunc:      mocks.NewBookStorage,
			ServiceNames: Services{"db"},
		},
		"book_storage_statistics_provider": Node{
			Ev: Event{
				Name:    "add_stats_provider",
				Service: "book_storage",
			},
		},
		"authors_storage": Node{
			NewFunc:      mocks.NewAuthorsStorage,
			ServiceNames: Services{"db"},
		},
		"authors_storage_statistics_provider": Node{
			Ev: Event{
				Name:    "add_stats_provider",
				Service: "authors_storage",
			},
		},
		"stats_provide_observer": Node{
			Ob: Observer{
				Event: "statistics_provider",
				Name:  "statistics_gateway",
				Callback: func(sg *mocks.StatisticsGateway, sp mocks.StatisticsProvider) {
					sg.AddStatisticsProvider(sp)
				},
			},
		},
		"static_files_url": Node{
			Constr: func(c Container) (interface{}, error) {
				var config mocks.Config
				c.Scan("Tree", &config)

				return config.GetValue("staticFilesUrl"), nil
			},
		},
		"book_link_provider": Node{
			NewFunc:      mocks.NewBookLinkProvider,
			ServiceNames: Services{"static_files_url", "book_finder_declared_statically"},
		},
		"web_fetcher":     Node{NewFunc: mocks.NewWebFetcher},
		"in_memory_cache": Node{NewFunc: mocks.NewInMemoryCache},
	}
}
