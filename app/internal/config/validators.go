package config

import (
	"fmt"
	"net/url"
)

func ValidateAdaptiveConfiguration(c AdaptiveConfig) error {
    if c.MaxExpectedConnections <= 0 {
        return fmt.Errorf("max expected connections must be positive")
    }
    if c.LatencyTarget <= 0 {
        return fmt.Errorf("latency target must be positive")
    }
    if c.ActiveWeight < 0 || c.LatencyWeight < 0 {
        return fmt.Errorf("adaptive weights must not be negative")
    }
    if c.ActiveWeight+c.LatencyWeight == 0 {
        return fmt.Errorf("at least one adaptive weight must be positive")
    }
    return nil
}

func ParseBackendURL(raw string) (*url.URL, error) {
	target, err := url.Parse(raw)
	if err != nil {
		return nil, fmt.Errorf("parse backend URL %q: %w", raw, err)
	}

	if target.Scheme != "http" && target.Scheme != "https" {
		return nil, fmt.Errorf("backend URL %q must use http or https", raw)
	}

	if target.Host == "" {
		return nil, fmt.Errorf("backend URL %q must include a host", raw)
	}

	return target, nil
}

func ParseAlgorithm(value string) (string, error) {
    switch value {
    case AlgorithmRoundRobin,
        AlgorithmHealthyRoundRobin,
        AlgorithmLeastConnections,
        AlgorithmLatencyAware,
        AdaptiveSelection:
        return value, nil
    case "":
        return AdaptiveSelection, nil
    default:
        return "", fmt.Errorf("unsupported algorithm: %q", value)
    }
}