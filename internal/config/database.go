package config

type _Database struct {
	Type     string `yaml:"type"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	SSLMode  string `yaml:"ssl_mode"`
	DBFile   string `yaml:"db_file"`
}

var Database = Load("database", &_Database{
	Type:   "sqlite",
	DBFile: "server.db",
})
