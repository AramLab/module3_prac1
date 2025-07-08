package service

type TaskRequest struct {
	Title       string `json:"title" validate:"required,min=3,max=100"`
	Description string `json:"description" validate:"required,min=5,max=500"`
	Status      string `json:"status" validate:"required,oneof=pending in_progress completed"`
}

type TaskRequestById struct {
	ID int `json:"id" validate:"required,gt=0"`
}

type TaskRequestUpdate struct {
	ID          int    `json:"id" validate:"required,gt=0"`
	Title       string `json:"title" validate:"omitempty,min=3,max=100"`
	Description string `json:"description" validate:"omitempty,min=5,max=500"`
	Status      string `json:"status" validate:"omitempty,oneof=pending in_progress completed"`
}
