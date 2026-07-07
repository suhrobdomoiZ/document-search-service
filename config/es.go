package config

import "strconv"

const DefaultSearchSize int = 10

type EsConfig struct {
	address    string
	searchSize int
}

func NewESConfig() *EsConfig {
	port := EsPort.MustGet()
	searchSizeString := EsSearchSize.MustGet()

	searchSize, err := strconv.Atoi(searchSizeString)
	if err != nil {
		searchSize = DefaultSearchSize
	}

	return &EsConfig{address: "http://localhost:" + port, searchSize: searchSize}
}

func (c *EsConfig) Address() string {
	return c.address
}

func (c *EsConfig) SearchSize() int {
	return c.searchSize
}
