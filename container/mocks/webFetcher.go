package mocks

//WebFetcher simulates downloading of a data from a url
type WebFetcher struct{}

func NewWebFetcherPtr() *WebFetcher {
	return &WebFetcher{}
}

func NewWebFetcher() WebFetcher {
	return WebFetcher{}
}

func NewWrongWebFetcherAsPtr() *BookCreator {
	return &BookCreator{}
}

func NewWrongWebFetcher() BookCreator {
	return BookCreator{}
}

func NewWrongWebFetcherAsSlice() []string {
	return []string{}
}

func NewWrongWebFetcherAsChan() chan bool {
	return make(chan bool)
}

func NewWrongWebFetcherAsMap() map[int]bool {
	return make(map[int]bool)
}

func NewWrongWebFetcherAsInterface() interface{} {
	var a interface{} = 1
	return a
}

//Fetch fake method to download something from a url
func (wf *WebFetcher) Fetch(url string) string {
	return "Fetched from " + url
}
