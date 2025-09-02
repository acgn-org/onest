package config

import (
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type _Server struct {
	Host        string `yaml:"host"`
	Port        uint16 `yaml:"port"`
	Timeout     uint   `yaml:"timeout"`
	LogLevel    string `yaml:"log_level"`
	LogRingSize int    `yaml:"log_ring_size"`
	FilePerm    string `yaml:"file_perm"`
}

var Server = LoadScoped("server", &_Server{
	Host:        "0.0.0.0",
	Port:        80,
	Timeout:     30,
	LogLevel:    "info",
	LogRingSize: 500,
	FilePerm:    "0777",
})

var FilePerm os.FileMode

func init() {
	lv, err := log.ParseLevel(Server.Get().LogLevel)
	if err != nil {
		log.Fatalln("failed to parse log level:", err)
	}
	log.SetLevel(lv)

	perm, err := strconv.ParseUint(Server.Get().FilePerm, 8, 32)
	if err != nil {
		log.Fatalf("invalid file perm '%s'", Server.Get().FilePerm)
	}
	FilePerm = os.FileMode(perm)
}
