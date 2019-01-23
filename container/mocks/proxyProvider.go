package mocks

type ProxyProvider interface {
	GetProxy() string
}

type ProxyRotator struct {
	proxyProviders []ProxyProvider
}

func NewProxyRotator(proxyProviders ...ProxyProvider) ProxyRotator {
	return ProxyRotator{proxyProviders}
}

type HardcodedProxyProvider struct{}

func NewHardcodedProxyProvider() HardcodedProxyProvider {
	return HardcodedProxyProvider{}
}

func (hpp HardcodedProxyProvider) GetProxy() string {
	return "someProxy"
}
