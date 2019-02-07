package mocks

type WebfetcherCaller struct {
	webFetcher *WebFetcher
}

func NewWebfetcherCaller(webFetcher *WebFetcher) *WebfetcherCaller {
	return &WebfetcherCaller{webFetcher: webFetcher}
}
