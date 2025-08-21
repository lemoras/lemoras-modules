package initialize

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type Bucket struct {
	gorm.Model
	BucketId   uuid.UUID `json:"bucketId"`
	UserId     uuid.UUID `json:"userId"`
	RoleId     int       `json:"roleId"`
	AppId      int       `json:"appId"`
	MerchantId uuid.UUID `json:"merchantId"`
	BucketName string    `json:"bucketName"`
	EmptySize  float64   `json:"emptySize"`
	SizeLimit  float64   `json:"sizeLimit"`
}

type BucketItem struct {
	gorm.Model
	BucketItemId uuid.UUID `json:"bucketItemId"`
	BucketId     uuid.UUID `json:"bucketId"`
	ParentId     uuid.UUID `json:"parentId"`
	ItemType     int       `json:"itemId"`
	ItemName     string    `json:"itemName"`
	ItemSize     float64   `json:"itemSize"`
	ItemURL      string    `json:"itemURL"`
}

type Note struct {
	gorm.Model
	NoteId     uuid.UUID `json:"noteId"`
	UserId     uuid.UUID `json:"userId"`
	RoleId     int       `json:"roleId"`
	AppId      int       `json:"appId"`
	MerchantId uuid.UUID `json:"merchantId"`
	Category   int       `json:"category"`
	TitleName  string    `json:"titleName"`
	DetailText string    `json:"detailText"`
}

type TokenRole struct {
	UserId     uuid.UUID
	RoleId     int
	MerchantId uuid.UUIDs
	AppId      int
}
