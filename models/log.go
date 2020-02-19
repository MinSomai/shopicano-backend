package models

import "time"

type Log struct {
	ID        string    `json:"id" gorm:"column:id;primary_key;not null"`
	Label     string    `json:"label" gorm:"column:label;index;not null"`
	Path      string    `json:"path" gorm:"column:path;index;not null"`
	Status    int       `json:"status" gorm:"column:status;index;not null"`
	Size      int64     `json:"size" gorm:"column:size;index;not null"`
	IP        string    `json:"ip" gorm:"column:ip;not null"`
	User      string    `json:"user" gorm:"user"`
	TimeTaken float64   `json:"time_taken" gorm:"column:time_taken;index;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;index;not null"`
}

func (l *Log) TableName() string {
	return "logs"
}
