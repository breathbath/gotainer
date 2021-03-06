package mocks

//BookDownloader simulates books downloading from a url
type BookDownloader struct {
	cache            Cache
	boolLinkProvider BookLinkProvider
	bookFinder       BookFinder
	downloadsCount   int
	webFetcher       *WebFetcher
}

//NewBookDownloader constructor for BookDownloader
func NewBookDownloader(cache Cache, boolLinkProvider BookLinkProvider, bookFinder BookFinder, webFetcher *WebFetcher) *BookDownloader {
	return &BookDownloader{cache, boolLinkProvider, bookFinder, 0, webFetcher}
}

//DownloadBook get a book by id
func (d *BookDownloader) DownloadBook(id string) string {
	book, _ := d.bookFinder.FindBook(id)
	link := d.boolLinkProvider.GetLink(id)

	d.cache.Cache(book)
	d.downloadsCount++

	return d.webFetcher.Fetch(link)
}
