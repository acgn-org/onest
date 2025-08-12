package config

import "github.com/acgn-org/onest/internal/logfield"

type _Telegram struct {
	ApiId               int32  `yaml:"api_id"`
	ApiHash             string `yaml:"api_hash"`
	DataFolder          string `yaml:"data_folder"`
	MaxParallelDownload uint8  `yaml:"max_parallel_download"`
	MaxDownloadError    uint32 `yaml:"max_download_error"`
	ScanThresholdDays   uint16 `yaml:"scan_threshold_days"`
}

var Telegram = LoadScoped("telegram", &_Telegram{
	DataFolder:          "tdlib",
	MaxParallelDownload: 1,
	MaxDownloadError:    5,
	ScanThresholdDays:   32,
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
