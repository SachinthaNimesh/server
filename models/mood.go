package models

import (
	"time"
)

type Mood struct {
	ID         int       `json:"id" gorm:"primaryKey"`
	StudentID  int       `json:"student_id" gorm:"index"`
	RecordedAt time.Time `json:"recorded_at" gorm:"index;default:CURRENT_TIMESTAMP"`
	Emotion    string    `json:"emotion"`
	IsDaily    bool      `json:"is_daily"`
}

func (Mood) TableName() string {
	return "mood"
}
