package examples

type WebFetcher struct{}

func NewWebFetcher() *WebFetcher {
	return &WebFetcher{}
}

func (wf *WebFetcher) Fetch(url string) string {
	return "Fetched from " + url
}
