package pgclient

type Config struct {
	DSN             string `yaml:"dsn"`
	MaxIdleConns    int    `yaml:"max_idle_conns"`
	ConnMaxLifetime int64  `yaml:"conn_max_lifetime"`
	ConnMaxIdleTime int64  `yaml:"conn_max_idle_time"`
	MaxOpenConns    int    `yaml:"max_open_conns"`
}
