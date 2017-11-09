package examples

import "github.com/breathbath/gotainer/container"

func CreateContainer() container.Container {
	runtimeContainer := container.NewRuntimeContainer()

	runtimeContainer.AddConstructor("book_creator", func(c container.Container) (interface{}, error) {
		return BookCreator{}, nil
	})

	runtimeContainer.AddNewMethod("config", NewConfig)

	runtimeContainer.AddConstructor("db",  func(c container.Container) (interface{}, error){
		var config Config
		c.Scan("config", &config)
		connectionString := config.GetValue("fakeDbConnectionString")

		return NewFakeDb(connectionString), nil
	})

	runtimeContainer.AddConstructor("book_storage", func(c container.Container) (interface{}, error) {
		var db FakeDb
		c.Scan("db", &db)
		return NewBookStorage(db), nil
	})

	runtimeContainer.AddConstructor("book_finder", func(c container.Container) (interface{}, error) {
		var bc BookCreator
		c.Scan("book_creator", &bc)

		var bs BookStorage
		c.Scan("book_storage", &bs)

		return NewBookFinder(bs, bc), nil
	})

	runtimeContainer.AddNewMethod("book_finder_declared_statically", NewBookFinder, "book_storage", "book_creator")

	runtimeContainer.AddConstructor("static_files_url", func(c container.Container) (interface{}, error) {
		var config Config
		c.Scan("config", &config)

		return config.GetValue("staticFilesUrl"), nil
	})

	runtimeContainer.AddNewMethod("book_link_provider", NewBookLinkProvider, "static_files_url", "book_finder_declared_statically")

	runtimeContainer.AddNewMethod("web_fetcher", NewWebFetcher)
	runtimeContainer.AddNewMethod("in_memory_cache", NewInMemoryCache)

	runtimeContainer.AddNewMethod("book_downloader", NewBookDownloader, "in_memory_cache", "book_link_provider", "book_finder", "web_fetcher")


	runtimeContainer.AddConstructor("wrong_book_creator", func(c container.Container) (interface{}, error) {
		return 123, nil
	})

	runtimeContainer.AddNewMethod("book_shelve", NewBookShelve)

	return runtimeContainer
}
