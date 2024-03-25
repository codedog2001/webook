//go:build !k8s

package config

var Config = config{
	DB: DBConfig{
		DSN: "localhost:13316",
	},
	Redis: RedisConfig{
		Add: "localhost:6379",
	},
}
