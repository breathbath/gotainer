package container

import (
	"github.com/breathbath/gotainer/container/mocks"
)

//CreateContainer gives a container example
func CreateContainer() *RuntimeContainer {
	runtimeContainer := NewRuntimeContainer()

	runtimeContainer.AddConstructor("book_creator", func(c Container) (interface{}, error) {
		return mocks.BookCreator{}, nil
	})

	runtimeContainer.AddNewMethod("config", mocks.NewConfig)

	runtimeContainer.AddConstructor("db", func(c Container) (interface{}, error) {
		var config mocks.Config
		c.Scan("config", &config)
		connectionString := config.GetValue("fakeDbConnectionString")

		return mocks.NewFakeDb(connectionString), nil
	})

	runtimeContainer.AddConstructor("book_storage", func(c Container) (interface{}, error) {
		db := c.Get("db", true).(*mocks.FakeDb)
		return mocks.NewBookStorage(db), nil
	})
	runtimeContainer.RegisterDependencyEvent("statistics_provider", "book_storage")

	runtimeContainer.AddConstructor("book_finder", func(c Container) (interface{}, error) {
		var bc mocks.BookCreator
		c.Scan("book_creator", &bc)

		var bs mocks.BookStorage
		c.Scan("book_storage", &bs)

		return mocks.NewBookFinder(bs, bc), nil
	})

	runtimeContainer.AddNewMethod("book_finder_declared_statically", mocks.NewBookFinder, "book_storage", "book_creator")

	runtimeContainer.AddConstructor("static_files_url", func(c Container) (interface{}, error) {
		var config mocks.Config
		c.Scan("config", &config)

		return config.GetValue("staticFilesUrl"), nil
	})

	runtimeContainer.AddNewMethod("book_link_provider", mocks.NewBookLinkProvider, "static_files_url", "book_finder_declared_statically")

	runtimeContainer.AddNewMethod("web_fetcher", mocks.NewWebFetcher)
	runtimeContainer.AddNewMethod("in_memory_cache", mocks.NewInMemoryCache)

	runtimeContainer.AddNewMethod("book_downloader", mocks.NewBookDownloader, "in_memory_cache", "book_link_provider", "book_finder", "web_fetcher")

	runtimeContainer.AddConstructor("wrong_book_creator", func(c Container) (interface{}, error) {
		return 123, nil
	})

	runtimeContainer.AddNewMethod("book_shelve", mocks.NewBookShelve)

	runtimeContainer.AddNewMethod("authors_storage", mocks.NewAuthorsStorage, "db")
	runtimeContainer.RegisterDependencyEvent("statistics_provider", "authors_storage")

	runtimeContainer.AddNewMethod("statistics_gateway", mocks.NewStatisticsGateway)

	runtimeContainer.AddDependencyObserver("statistics_provider", "statistics_gateway", func(sg *mocks.StatisticsGateway, sp mocks.StatisticsProvider) {
		sg.AddStatisticsProvider(sp)
	})

	return runtimeContainer
}
