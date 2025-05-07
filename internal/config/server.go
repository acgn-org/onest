package config

import (
	nested "github.com/antonfisher/nested-logrus-formatter"
	log "github.com/sirupsen/logrus"
)

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

func init() {
	lv, err := log.ParseLevel(Server.LogLevel)
	if err != nil {
		log.Fatalln("failed to parse log level:", err)
	}
	log.SetLevel(lv)
	log.SetFormatter(&nested.Formatter{
		TimestampFormat: "01/02 15:04:05",
	})
}
