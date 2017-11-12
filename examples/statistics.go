package examples

//StatisticsProvider gives a statistic metric
type StatisticsProvider interface {
	GetStatistics() (string, int)
}

//StatisticsGateway to simulate the main statistics engine
type StatisticsGateway struct {
	statisticsProviders []StatisticsProvider
}

//Constructor for StatisticsGateway
func NewStatisticsGateway() *StatisticsGateway {
	return &StatisticsGateway{statisticsProviders: []StatisticsProvider{}}
}

//AddStatisticsProvider adds a StatisticsProvider to the StatisticsGateway
func (sg *StatisticsGateway) AddStatisticsProvider(sp StatisticsProvider) {
	sg.statisticsProviders = append(sg.statisticsProviders, sp)
}

//CollectStatistics calls all StatisticsProvider to collect metrics
func (sg *StatisticsGateway) CollectStatistics() map[string]int {
	result := make(map[string]int)

	for _, sp := range sg.statisticsProviders {
		statName, statCount := sp.GetStatistics()
		result[statName] = statCount
	}

	return result
}
