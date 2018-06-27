package models

import (
	"time"
)

type User struct {
	DynamicPassTime time.Time//动态密码过期时间
	DynamicPass string//动态密码
	UserID string
	Service string
	NickName string
}
type IMRecord struct {
	Id int64`gorm:"Column:id" json:"id"`
	Content string`gorm:"Column:content" json:"content"`
	Sender string`gorm:"Column:sender" json:"sender"`
	Recver string`gorm:"Column:recver" json:"recver"`
	Time string`gorm:"Column:time" json:"time"`
	Avatar string`gorm:"Column:avatar" json:"avatar"`
	Sendername string`gorm:"Column:sendername" json:"sendername"`
	Recvername string`gorm:"Column:recvername" json:"recvername"`
	Offline int8`gorm:"Column:offline" json:"offline"`
}
func (a *IMRecord)TableName()string  {
	return "IMRecord"
}
type OfflineMsg struct {
	Module 		string
	Method 		string
	Content        string
	Sender     string
	Recver     string
	Time       string
	Avatar     string
	Sendername string
	Recvername string
}

type ChatLog struct {
	Total int`json:"total"`
	Size int`json:"size"`
	Records []IMRecord`json:"records"`
}
type MessageResult struct {
	Module string
	Method string
	Service string
	Message string
}