package examples

type WebFetcher struct{}

//Simulates downloading of a data from a url
func NewWebFetcher() *WebFetcher {
	return &WebFetcher{}
}

//Fake method to download something from a url
func (wf *WebFetcher) Fetch(url string) string {
	return "Fetched from " + url
}
