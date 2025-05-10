package repository

type Item struct {
	ID        uint  `gorm:"primarykey"`
	RemoteID  uint  `gorm:"not null"`
	ChannelID int64 `gorm:"not null"`

	Name   string `gorm:"not null"`
	Regexp string `gorm:"not null"`

	DateStart int32 `gorm:"not null"`
	DateEnd   int32 `gorm:"not null"`

	Process int64 `gorm:"not null"`

	Priority   int32  `gorm:"not null"`
	TargetPath string `gorm:"not null"`
}

type ItemRepository struct {
	Repository
}
