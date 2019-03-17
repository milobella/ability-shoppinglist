package config

type ServerConfiguration struct {
	Port int 		`env:"SERVER_PORT" default:"8080"`
	LogLevel string `id:"log_level" env:"SERVER_LOG_LEVEL" default:"<root>=ERROR"`
}

type ToolConfiguration struct {
	Host string
	Port int
}