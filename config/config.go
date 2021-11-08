package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	RestPort int    `envconfig:"rest_port" default:"80"`
	TimeOut  int    `envconfig:"time_out" default:"3"`
	PassKey  string `envconfig:"pass_key" default:"ThMNCpmKpuCcx4XWFCBEzsVG6rnKFquw"`

	WaktuSholatHost string `envconfig:"waktu_sholat_host" default:"https://api.pray.zone"`
}

func Get() Config {
	cfg := Config{}

	envconfig.MustProcess("", &cfg)

	return cfg
}
