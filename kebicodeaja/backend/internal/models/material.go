package models

import "time"

type Material struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ClassID   uint      `gorm:"not null;index" json:"class_id"`
	Title     string    `gorm:"size:200;not null" json:"title"`
	Content   string    `gorm:"type:longtext" json:"content"`
	CreatedAt time.Time `json:"created_at"`

	Class       Class                `gorm:"foreignKey:ClassID;constraint:OnDelete:CASCADE" json:"-"`
	Attachments []MaterialAttachment `gorm:"foreignKey:MaterialID;constraint:OnDelete:CASCADE" json:"attachments,omitempty"`
}

type MaterialAttachment struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	MaterialID     uint           `gorm:"not null;index" json:"material_id"`
	AttachmentType AttachmentType `gorm:"type:enum('video','ebook','file','link');not null" json:"attachment_type"`
	Title          string         `gorm:"size:200" json:"title"`
	URL            string         `gorm:"size:500;not null" json:"url"`
	FileSize       int64          `json:"file_size"`
	CreatedAt      time.Time      `json:"created_at"`
}
