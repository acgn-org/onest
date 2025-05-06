package realsearch

type TimeMachineItem struct {
	ID     uint `json:"id"`
	RuleID uint `json:"rule_id"`

	NameCN       string `json:"name"`
	NameCNParsed string `json:"name_cn"`
	NameEN       string `json:"name_en"`

	DateStart int32 `json:"date_start"`
	DateEnd   int32 `json:"date_end"`
}

type ScheduleV2 struct {
	TimeMachineItem
	Status string `json:"status"`
}

type RawInfo struct {
	ID                uint   `json:"id"`
	ItemID            uint   `json:"item_id,omitempty"`
	ChannelID         int64  `json:"channel_id"`
	ChannelName       string `json:"channel_name"`
	Size              int32  `json:"size"`
	Text              string `json:"text"`
	FileSuffix        string `json:"file_suffix"`
	MsgID             int64  `json:"msg_id"`
	SupportsStreaming bool   `json:"supports_streaming"`
	Link              string `json:"link"`
	Date              int32  `json:"date"`
}
