package examples

//Simulates downloading of a data from a url
type WebFetcher struct{}

//NewWebFetcher constructor for WebFetcher
func NewWebFetcher() *WebFetcher {
	return &WebFetcher{}
}

//Fake method to download something from a url
func (wf *WebFetcher) Fetch(url string) string {
	return "Fetched from " + url
}
