package container

import (
	"github.com/breathbath/gotainer/container/mocks"
	"testing"
)

func TestDynamicDependencies(t *testing.T) {
	exampleConfig = Tree{
		Node{
			Parameters: map[string]interface{}{
				"url1": "some_url1",
				"url2": "some_url2",
				"domain": "some_domain",
			},
		},
		Node{
			ID:           "urlProvider",
			NewFunc:      mocks.NewUrlProvider,
			ServiceNames: Services{"url1", "url2"},
		},
		Node{
			ID:           "urlProviderWithDomain",
			NewFunc:      mocks.NewUrlProviderWithDomain,
			ServiceNames: Services{"domain", "url1", "url2"},
		},
	}

	cont := RuntimeContainerBuilder{}.BuildContainerFromConfig(exampleConfig)
	urlsProvider := cont.Get("urlProvider", true).(mocks.UrlProvider)
	urls := urlsProvider.GetUrls()
	if urls[0] != "some_url1" || urls[1] != "some_url2" {
		t.Errorf("Unexpected urls output %v for dynamic dependencies list", urls)
	}

	urlsProviderWithDomain := cont.Get("urlProviderWithDomain", true).(mocks.UrlProvider)
	urls = urlsProvider.GetUrls()
	if urls[0] != "some_url1" || urls[1] != "some_url2" {
		t.Errorf("Unexpected urls output %v for variadic dependencies list", urls)
	}

	if urlsProviderWithDomain.GetDomain() != "some_domain" {
		t.Errorf("Unexpected domain output %s for variadic dependencies list", urlsProviderWithDomain.GetDomain())
	}
}
