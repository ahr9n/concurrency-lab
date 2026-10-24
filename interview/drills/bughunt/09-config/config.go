//go:build bughunt

package config

import "runtime"

type Config struct {
	Addr    string
	Timeout int
}

var (
	current *Config
	ready   bool
)

func Load() {
	current = &Config{Addr: "db:5432", Timeout: 30}
	ready = true
}

func Get() *Config {
	for !ready {
		runtime.Gosched()
	}
	return current
}
