package config

type _Telegram struct {
	ApiId               int32  `yaml:"api_id"`
	ApiHash             string `yaml:"api_hash"`
	DataFolder          string `yaml:"data_folder"`
	MaxParallelDownload uint8  `yaml:"max_parallel_download"`
}

var Telegram = Load("telegram", &_Telegram{
	DataFolder:          "tdlib",
	MaxParallelDownload: 3,
})

func init() {
	if Telegram.MaxParallelDownload == 0 {
		Telegram.MaxParallelDownload = 1
	}
}
