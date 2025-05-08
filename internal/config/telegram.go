package config

type _Telegram struct {
	DataFolder string `yaml:"data_folder"`
	ApiId      int32  `yaml:"api_id"`
	ApiHash    string `yaml:"api_hash"`
}

var Telegram = Load("telegram", &_Telegram{
	DataFolder: "tdlib",
})
