package mocks

type WebfetcherCaller struct {
	webFetcherPtr *WebFetcher
	webFetcher    WebFetcher
}

func NewWebfetcherCallerByPtr(webFetcher *WebFetcher) *WebfetcherCaller {
	return &WebfetcherCaller{webFetcherPtr: webFetcher}
}

func NewWebfetcherCaller(webFetcher WebFetcher) *WebfetcherCaller {
	return &WebfetcherCaller{webFetcher: webFetcher}
}
