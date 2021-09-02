package model

type Base struct {
	Id
	Tail
}

type Id struct {
	ID int64 `json:"id" gorm:"primary_key"`
}

type Oid struct {
	OID int64 `json:"oid" gorm:"column:oid"`
}

type Uid struct {
	UID int64 `json:"uid" gorm:"column:uid"`
}

type Tail struct {
	Time int64 `json:"time"`
	//删除标识 1：删除 0：正常
	IsDelete int `json:"is_delete" gorm:"default:0"`
	//删除时间
	DeleteTime int64 `json:"delete_time"`
}
