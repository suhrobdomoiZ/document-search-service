package configkey

import (
	"fmt"
	"os"
)

type Key string

func (k Key) MustGet() string {
	val := os.Getenv(string(k))
	if val == "" {
		panic(fmt.Sprintf("config.MustGet: %s is required", string(k)))
	}

	return val
}

func (k Key) Get(def string) string {
	val := os.Getenv(string(k))
	if val == "" {
		val = def
	}

	return val
}
