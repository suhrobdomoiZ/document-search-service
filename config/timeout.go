package config

import (
	"strconv"
)

type TimeoutConfig struct {
	serverReadTimeout  int
	serverWriteTimeout int
	serverIdleTimeout  int
	shutdownCtxTimeout int
}

func NewTimeoutConfig() *TimeoutConfig {
	rawServerReadTimeout := ServerReadTimeout.Get("5")
	rawServerWriteTimeout := ServerWriteTimeout.Get("10")
	rawServerIdleTimeout := ServerIdleTimeout.Get("120")
	rawShutdownCtxTimeout := ShutdownCtxTimeout.Get("5")

	intServerReadTimeout, err := strconv.Atoi(rawServerReadTimeout)
	if err != nil {
		intServerReadTimeout = 5
	}

	intServerWriteTimeout, err := strconv.Atoi(rawServerWriteTimeout)
	if err != nil {
		intServerWriteTimeout = 10
	}

	intServerIdleTimeout, err := strconv.Atoi(rawServerIdleTimeout)
	if err != nil {
		intServerIdleTimeout = 120
	}

	intShutdownCtxTimeout, err := strconv.Atoi(rawShutdownCtxTimeout)
	if err != nil {
		intShutdownCtxTimeout = 5
	}

	return &TimeoutConfig{
		serverReadTimeout:  intServerReadTimeout,
		serverWriteTimeout: intServerWriteTimeout,
		serverIdleTimeout:  intServerIdleTimeout,
		shutdownCtxTimeout: intShutdownCtxTimeout,
	}
}

func (t *TimeoutConfig) ServerReadTimeout() int {
	return t.serverReadTimeout
}

func (t *TimeoutConfig) ServerWriteTimeout() int {
	return t.serverWriteTimeout
}

func (t *TimeoutConfig) ServerIdleTimeout() int {
	return t.serverIdleTimeout
}

func (t *TimeoutConfig) ShutdownCtxTimeout() int {
	return t.shutdownCtxTimeout
}
