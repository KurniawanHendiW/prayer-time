package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	// Service
	Host       string `envconfig:"HOST" default:""`
	PORT       int    `envconfig:"PORT" default:"80"`
	TimeOut    int    `envconfig:"TIME_OUT" default:"3"`
	ExpiredKey int64  `envconfig:"EXPIRED_KEY" default:"600"` // in seconds
	DebugLog   bool   `envconfig:"DEBUG_LOG" default:"true"`

	// Redis
	RedisTlsUrl    string `envconfig:"REDIS_TLS_URL" default:""` // format tls url, redis://:<password>@<host>:<port>
	RedisHost      string `envconfig:"REDIS_HOST" default:""`
	RedisPort      string `envconfig:"redis_port" default:""`
	RedisPassword  string `envconfig:"REDIS_PASSWORD" default:""`
	RedisTimeout   int    `envconfig:"REDIS_TIMEOUT" default:"3"`
	RedisMaxIdle   int    `envconfig:"REDIS_MAX_IDLE" default:"8"`
	RedisMaxActive int    `envconfig:"REDIS_MAX_ACTIVE" default:"10"`

	// third party
	WaktuSholatHost string `envconfig:"WAKTU_SHOLAT_HOST" default:"https://waktusholat.org"`
	ApiPrayZoneHost string `envconfig:"API_PRAY_ZONE_HOST" default:"https://api.pray.zone"`
}

func Get() Config {
	cfg := Config{}

	envconfig.MustProcess("", &cfg)

	return cfg
}
