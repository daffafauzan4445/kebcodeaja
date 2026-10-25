package dto

type CreateClassRequest struct {
	Name        string `json:"name" binding:"required,min=3,max=150"`
	Description string `json:"description" binding:"max=500"`
}

type JoinClassRequest struct {
	ClassCode string `json:"class_code" binding:"required"`
}
