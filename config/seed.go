package config

type SeedConfig struct {
	dataPath string
}

func NewSeedConfig() *SeedConfig {
	dataPath := DataPath.Get("./data/posts.csv")

	return &SeedConfig{dataPath}
}

func (c *SeedConfig) DataPath() string {
	return c.dataPath
}
