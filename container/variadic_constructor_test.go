package container

import (
	"github.com/breathbath/gotainer/container/mocks"
	"testing"
)

func TestVariadicConfigDependencies(t *testing.T) {
	exampleConfig = Tree{
		Node{
			Parameters: map[string]interface{}{
				"url1":   "some_url1",
				"url2":   "some_url2",
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
		Node{
			ID:           "urlProviderWithDomainWithoutUrls",
			NewFunc:      mocks.NewUrlProviderWithDomain,
			ServiceNames: Services{"domain"},
		},
	}

	cont, err := RuntimeContainerBuilder{}.BuildContainerFromConfig(exampleConfig)
	if err != nil {
		t.Error(err)
		return
	}

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

	urlProviderWithDomainWithoutUrls := cont.Get("urlProviderWithDomainWithoutUrls", true).(mocks.UrlProvider)
	if urlProviderWithDomainWithoutUrls.GetDomain() != "some_domain" {
		t.Errorf("Unexpected domain output %s for variadic dependencies list", urlsProviderWithDomain.GetDomain())
	}
	urls = urlProviderWithDomainWithoutUrls.GetUrls()
	if len(urls) != 0 {
		t.Errorf("Expected zero urls length but %v urls list are provided", urls)
	}
}

func TestWrongVariadicConfigDependencies(t *testing.T) {
	exampleConfig = Tree{
		Node{
			Parameters: map[string]interface{}{
				"count": 2,
			},
		},
		Node{
			ID:      "bookShelve",
			NewFunc: mocks.NewBookShelve,
		},
		Node{
			ID:           "urlProvider",
			NewFunc:      mocks.NewUrlProvider,
			ServiceNames: Services{"count", "bookShelve"},
		},
	}

	defer ExpectPanic(
		t,
		"Cannot use the provided dependency 'count' of type 'int' as 'string' in the Constr function call [check 'urlProvider' service];\n"+
		"Cannot use the provided dependency 'bookShelve' of type '*mocks.BookShelve' as 'string' in the Constr function call [check 'urlProvider' service]",
	)

	cont, err := RuntimeContainerBuilder{}.BuildContainerFromConfig(exampleConfig)
	if err != nil {
		t.Error(err)
		return
	}
	cont.Get("urlProvider", true)
}

func TestPartiallyWrongVariadicConfigDependencies(t *testing.T) {
	exampleConfig = Tree{
		Node{
			Parameters: map[string]interface{}{
				"domain": "my_domain",
				"some_list": []string{"listItemA", "listItemB"},
			},
		},
		Node{
			ID:           "urlProvider",
			NewFunc:      mocks.NewUrlProviderWithDomain,
			ServiceNames: Services{"domain", "some_list"},
		},
	}

	defer ExpectPanic(
		t,
		"Cannot use the provided dependency 'some_list' of type '[]string' as 'string' in the Constr function call [check 'urlProvider' service]",
	)

	cont, err := RuntimeContainerBuilder{}.BuildContainerFromConfig(exampleConfig)
	if err != nil {
		t.Error(err)
		return
	}

	cont.Get("urlProvider", true)
}

func TestVariadicConstructor(t *testing.T) {
	cont := NewRuntimeContainer()
	stringParams := map[string]string{
		"url1": "MyUrl1",
		"url2": "MyUrl2",
	}
	err := RegisterParameters(cont, stringParams)
	if err != nil {
		t.Error(err)
	}

	cont.AddNewMethod("urlProvider", mocks.NewUrlProvider, "url1", "url2")

	urlsProvider := cont.Get("urlProvider", true).(mocks.UrlProvider)
	urls := urlsProvider.GetUrls()
	if urls[0] != "MyUrl1" || urls[1] != "MyUrl2" {
		t.Errorf("Unexpected urls output %v for dynamic dependencies list", urls)
	}
}
