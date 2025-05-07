package config

type _Server struct {
	Host     string `yaml:"host"`
	Port     uint16 `yaml:"port"`
	LogLevel string `yaml:"logLevel"`
}

var Server = Load("server", &_Server{
	Host:     "0.0.0.0",
	Port:     80,
	LogLevel: "info",
})
