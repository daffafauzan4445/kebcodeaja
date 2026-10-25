package models

import "time"

type Class struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	TeacherID   uint      `gorm:"not null;index" json:"teacher_id"`
	Name        string    `gorm:"size:150;not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	ClassCode   string    `gorm:"size:12;uniqueIndex;not null" json:"class_code"`
	CreatedAt   time.Time `json:"created_at"`

	Teacher User `gorm:"foreignKey:TeacherID;constraint:OnDelete:RESTRICT" json:"-"`
}

type ClassMember struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ClassID   uint      `gorm:"not null;uniqueIndex:idx_class_student" json:"class_id"`
	StudentID uint      `gorm:"not null;uniqueIndex:idx_class_student" json:"student_id"`
	JoinedAt  time.Time `gorm:"autoCreateTime" json:"joined_at"`

	Class   Class `gorm:"foreignKey:ClassID;constraint:OnDelete:CASCADE" json:"-"`
	Student User  `gorm:"foreignKey:StudentID;constraint:OnDelete:CASCADE" json:"-"`
}
