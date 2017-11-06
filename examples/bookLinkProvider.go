package examples

type BookLinkProvider struct {
	downloadUrl string
	bookFinder  BookFinder
}

func NewBookLinkProvider(downloadUrl string, bookFinder BookFinder) BookLinkProvider {
	return BookLinkProvider{downloadUrl, bookFinder}
}

func (blp BookLinkProvider) GetLink(id string) string {
	book, _ := blp.bookFinder.FindBook(id)
	return blp.downloadUrl + book.Title
}
