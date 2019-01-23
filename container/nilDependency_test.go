package container

import (
	"github.com/breathbath/gotainer/container/mocks"
	"testing"
)

func TestNilDependencyForPointerOrInterfaceConstructorArgument(t *testing.T) {
	cont := NewRuntimeContainer()
	cont.AddConstructor("nilDb", func(c Container) (i interface{}, e error) {
		return nil, nil
	})
	err := RegisterParameters(cont, map[string]interface{}{"nilCache": nil})
	if err != nil {
		t.Error(err)
	}

	cont.AddNewMethod("bookStorage", mocks.NewBookStorage, "nilDb")
	cont.AddNewMethod("cacheManager", mocks.NewCacheManager, "nilCache")

	cont.Get("bookStorage", true)
	cont.Get("cacheManager", true)
}

func TestNilDependencyForVariadicConstructors(t *testing.T) {
	cont, err := initContainerWithProxyDependencies(true)
	if err != nil {
		t.Error(err)
	}

	cont.Get("proxyRotatorFromList", true)
	cont.Get("proxyRotatorFromSingleProxyProvider", true)

	cont, err = initContainerWithProxyDependencies(false)
	if err != nil {
		t.Error(err)
	}
	cont.Get("proxyRotatorFromList", true)
	cont.Get("proxyRotatorFromSingleProxyProvider", true)
}

func initContainerWithProxyDependencies(isProxyEnabled bool) (*RuntimeContainer, error) {
	cont := NewRuntimeContainer()

	err := RegisterParameters(cont, map[string]bool{"isProxyEnabled": isProxyEnabled})
	if err != nil {
		return cont, err
	}

	cont.AddNewMethod("hardcodedProxyProvider", mocks.NewHardcodedProxyProvider)

	cont.AddConstructor("proxyProvidersList", func(c Container) (i interface{}, e error) {
		isProxyEnabled := c.Get("isProxyEnabled", true).(bool)
		if isProxyEnabled {
			hardcodedProxy := c.Get("hardcodedProxyProvider", true).(mocks.HardcodedProxyProvider)
			return []mocks.ProxyProvider{hardcodedProxy}, nil
		}

		//actually this should be []mocks.ProxyProvider{} but user might do this
		return []mocks.ProxyProvider{nil}, nil
	})

	cont.AddConstructor("mainProxyProvider", func(c Container) (i interface{}, e error) {
		isProxyEnabled := c.Get("isProxyEnabled", true).(bool)
		if isProxyEnabled {
			hardcodedProxy := c.Get("hardcodedProxyProvider", true).(mocks.HardcodedProxyProvider)
			return hardcodedProxy, nil
		}

		//actually it would be better to use NullObject implementation for ProxyProvider interface rather than
		//return nil, but user of library might do this unfortunate design
		return nil, nil
	})

	cont.AddConstructor("proxyRotatorFromList", func(c Container) (i interface{}, e error) {
		proxyProviderList := c.Get("proxyProvidersList", true).([]mocks.ProxyProvider)
		proxyRotator := mocks.NewProxyRotator(proxyProviderList...)
		return proxyRotator, nil
	})

	cont.AddNewMethod("proxyRotatorFromSingleProxyProvider", mocks.NewProxyRotator, "mainProxyProvider")
	return cont, nil
}
