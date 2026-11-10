package backend

type Backend struct {
	URL string
}

func Targets(targetURLs []string) []Backend {
	backends := make([]Backend, 0, len(targetURLs))

	for _, targetURL := range targetURLs {
		backends = append(backends, Backend{URL: targetURL})
	}

	return backends
}
