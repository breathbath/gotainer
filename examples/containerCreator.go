package examples

import "github.com/breathbath/gotainer/container"

func CreateContainer() container.Container {
	runtimeContainer := container.NewRuntimeContainer()

	runtimeContainer.AddNoArgumentsConstructor("book_creator", func() (interface{}, error) {
		return BookCreator{}, nil
	})

	runtimeContainer.AddTypedConstructor("config", NewConfig)

	runtimeContainer.AddConstructor("db",  func(c container.Container) (interface{}, error){
		var config Config
		c.GetTypedService("config", &config)
		connectionString := config.GetValue("fakeDbConnectionString")

		return NewFakeDb(connectionString), nil
	})

	runtimeContainer.AddConstructor("book_storage", func(c container.Container) (interface{}, error) {
		var db FakeDb
		c.GetTypedService("db", &db)
		return NewBookStorage(db), nil
	})

	runtimeContainer.AddConstructor("book_finder", func(c container.Container) (interface{}, error) {
		var bc BookCreator
		c.GetTypedService("book_creator", &bc)

		var bs BookStorage
		c.GetTypedService("book_storage", &bs)

		return NewBookFinder(bs, bc), nil
	})

	runtimeContainer.AddTypedConstructor("book_finder_declared_statically", NewBookFinder, "book_storage", "book_creator")

	runtimeContainer.AddConstructor("static_files_url", func(c container.Container) (interface{}, error) {
		var config Config
		c.GetTypedService("config", &config)

		return config.GetValue("staticFilesUrl"), nil
	})

	runtimeContainer.AddTypedConstructor("book_link_provider", NewBookLinkProvider, "static_files_url", "book_finder_declared_statically")

	runtimeContainer.AddTypedConstructor("web_fetcher", NewWebFetcher)
	runtimeContainer.AddTypedConstructor("in_memory_cache", NewInMemoryCache)

	runtimeContainer.AddTypedConstructor("book_downloader", NewBookDownloader, "in_memory_cache", "book_link_provider", "book_finder", "web_fetcher")


	runtimeContainer.AddNoArgumentsConstructor("wrong_book_creator", func() (interface{}, error) {
		return 123, nil
	})

	runtimeContainer.AddConstructor("wrong_book_finder", func(c container.Container) (interface{}, error) {
		var bc BookCreator
		c.GetTypedService("wrong_book_creator", &bc)

		var bs BookStorage
		c.GetTypedService("book_storage", &bs)

		return NewBookFinder(bs, bc), nil
	})

	return runtimeContainer
}
