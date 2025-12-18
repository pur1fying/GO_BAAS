package config

type DatabaseConfig struct {
	Host            string `yaml:"host"`
	Port            int    `yaml:"port"`
	Username        string `yaml:"username"`
	Password        string `yaml:"password"`
	Name            string `yaml:"name"`
	MaxOpenConn     int    `yaml:"max_open_conn"`
	MaxIdleConn     int    `yaml:"max_idle_conn"`
	ConnMaxLifetime int    `yaml:"conn_max_lifetime"`
}

func DefaultDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Host:            "127.0.0.1",
		Port:            3306,
		Username:        "",
		Password:        "",
		Name:            "GO_BAAS",
		MaxOpenConn:     100,
		MaxIdleConn:     10,
		ConnMaxLifetime: 60,
	}
}
