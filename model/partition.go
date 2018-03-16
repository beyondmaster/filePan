package model

import "time"

// partition 分区
type Partition struct {
	ID          uint      `gorm:"primary_key" json:"id"`
	Name        string    `json:"name"`
	CreatedBy   uint      `json:"createdBy"`
	CreatedTime time.Time `json:"createdTime"`
	IsTop       int       `json:"isTop"`
}

const MaxPartitionNameLen = 10
