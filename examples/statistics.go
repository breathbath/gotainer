package examples

type StatisticsProvider interface {
	GetStatistics() (string, int)
}

type StatisticsGateway struct {
	statisticsProviders []StatisticsProvider
}

func NewStatisticsGateway() *StatisticsGateway {
	return &StatisticsGateway{statisticsProviders: []StatisticsProvider{}}
}

func (sg *StatisticsGateway) AddStatisticsProvider(sp StatisticsProvider) {
	sg.statisticsProviders = append(sg.statisticsProviders, sp)
}

func (sg *StatisticsGateway) CollectStatistics() map[string]int {
	result := make(map[string]int)

	for _, sp := range sg.statisticsProviders {
		statName, statCount := sp.GetStatistics()
		result[statName] = statCount
	}

	return result
}
