package repository

type Download struct {
	ID uint `gorm:"primarykey"`

	ItemID uint `gorm:"index;not null"`
	Item   Item `gorm:"foreignKey:ItemID;constraint:OnDelete:CASCADE"`

	MsgID int64 `gorm:"not null"`
	Text  string
	Size  int32 `gorm:"not null"`
	Date  int32 `gorm:"not null"`
}

type DownloadRepository Repository
