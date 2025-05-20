package config

import "github.com/acgn-org/onest/internal/logfield"

type _Telegram struct {
	ApiId               int32  `yaml:"api_id"`
	ApiHash             string `yaml:"api_hash"`
	DataFolder          string `yaml:"data_folder"`
	MaxParallelDownload uint8  `yaml:"max_parallel_download"`
	MaxDownloadError    uint32 `yaml:"max_download_error"`
}

var Telegram = LoadScoped("telegram", &_Telegram{
	DataFolder:          "tdlib",
	MaxParallelDownload: 3,
	MaxDownloadError:    5,
})

func init() {
	telegramConfig := Telegram.Get()

	if telegramConfig.MaxParallelDownload == 0 {
		telegramConfig.MaxParallelDownload = 1
		err := Telegram.Save(telegramConfig)
		if err != nil {
			logfield.New(logfield.ComConfig).Fatalln("save telegram config failed:", err)
		}
	}
}
