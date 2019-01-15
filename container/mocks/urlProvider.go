package mocks

type UrlProvider struct {
	domain string
	urls []string
}

func NewUrlProvider(urls ...string) UrlProvider {
	return UrlProvider{urls: urls, domain: ""}
}

func NewUrlProviderWithDomain(domain string, urls ...string) UrlProvider {
	return UrlProvider{urls: urls, domain: domain}
}

func (up UrlProvider) GetUrls() []string {
	return up.urls
}

func (up UrlProvider) GetDomain() string {
	return up.domain
}
