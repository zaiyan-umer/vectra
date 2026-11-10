package config

type Config struct {
	Algorithm  string
	TargetURLs []string
}

const AlgorithmRoundRobin = "round_robin"
const AlgorithmHealthyRoundRobin = "healthy_round_robin"
const AlgorithmLeastConnections = "least_connections"
