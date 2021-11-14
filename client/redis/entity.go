package redis

type (
	Reply struct {
		Result interface{}
		Error  error
	}

	RedisConfig struct {
		Host      string
		Port      string
		Timeout   int
		MaxIdle   int
		MaxActive int
		Password  string
	}
)
