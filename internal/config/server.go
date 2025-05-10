package config

import (
	log "github.com/sirupsen/logrus"
)

type _Server struct {
	Host     string `yaml:"host"`
	Port     uint16 `yaml:"port"`
	LogLevel string `yaml:"log_level"`
}

var Server = LoadScoped("server", &_Server{
	Host:     "0.0.0.0",
	Port:     80,
	LogLevel: "info",
})

func init() {
	lv, err := log.ParseLevel(Server.Get().LogLevel)
	if err != nil {
		log.Fatalln("failed to parse log level:", err)
	}
	log.SetLevel(lv)
}
