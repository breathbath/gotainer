package mocks

//BookLinkProvider creates book links
type BookLinkProvider struct {
	downloadUrl string
	bookFinder  BookFinder
}

//NewBookLinkProvider BookLinkProvider constructor
func NewBookLinkProvider(downloadUrl string, bookFinder BookFinder) BookLinkProvider {
	return BookLinkProvider{downloadUrl, bookFinder}
}

//GetLink returns a link to a book by its id
func (blp BookLinkProvider) GetLink(id string) string {
	book, _ := blp.bookFinder.FindBook(id)
	return blp.downloadUrl + book.Title
}
