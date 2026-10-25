package models

import (
	"time"

	"gorm.io/datatypes"
)

type Challenge struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	ClassID       uint           `gorm:"not null;index" json:"class_id"`
	MaterialID    *uint          `gorm:"index" json:"material_id"`
	Title         string         `gorm:"size:200;not null" json:"title"`
	Description   string         `gorm:"type:text" json:"description"`
	ChallengeType ChallengeType  `gorm:"type:enum('multiple_choice','ordering','fill_blank');not null" json:"challenge_type"`
	CTAspect      CTAspect       `gorm:"column:ct_aspect;type:enum('decomposition','pattern_recognition','abstraction','algorithm');not null" json:"ct_aspect"`
	Difficulty    Difficulty     `gorm:"type:enum('easy','medium','hard');default:'easy'" json:"difficulty"`
	Content       datatypes.JSON `gorm:"type:json" json:"content"`
	AnswerKey     datatypes.JSON `gorm:"type:json" json:"-"`
	Hint          string         `gorm:"type:text" json:"hint"`
	MaxScore      int            `gorm:"default:100" json:"max_score"`
	CreatedAt     time.Time      `json:"created_at"`

	Class    Class     `gorm:"foreignKey:ClassID;constraint:OnDelete:CASCADE" json:"-"`
	Material *Material `gorm:"foreignKey:MaterialID;constraint:OnDelete:SET NULL" json:"-"`
}

type Submission struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	ChallengeID uint           `gorm:"not null;uniqueIndex:idx_sub_unique" json:"challenge_id"`
	StudentID   uint           `gorm:"not null;uniqueIndex:idx_sub_unique" json:"student_id"`
	AttemptNo   int            `gorm:"not null;default:1;uniqueIndex:idx_sub_unique" json:"attempt_no"`
	Answer      datatypes.JSON `gorm:"type:json" json:"answer"`
	IsCorrect   bool           `json:"is_correct"`
	Score       int            `json:"score"`
	SubmittedAt time.Time      `gorm:"autoCreateTime" json:"submitted_at"`

	Challenge Challenge `gorm:"foreignKey:ChallengeID;constraint:OnDelete:CASCADE" json:"-"`
	Student   User      `gorm:"foreignKey:StudentID;constraint:OnDelete:CASCADE" json:"-"`
}
