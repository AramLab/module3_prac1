package service

type TaskRequestById struct {
	ID int `json:"id" validate:"required,gt=0"`
}

type TaskRequestUpdate struct {
	ID          int    `json:"id" validate:"required,gt=0"`
	Title       string `json:"title" validate:"omitempty,min=3,max=100"`
	Description string `json:"description" validate:"omitempty,min=5,max=500"`
	Status      string `json:"status" validate:"omitempty,oneof=pending in_progress completed"`
}

type UserRequest struct {
	Username string `json:"username" validate:"required,min=3"`
	Password string `json:"password" validate:"required,min=6"`
}

type TaskRequest struct {
	UserID      int    `json:"user_id" validate:"required"`
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
	Status      string `json:"status" validate:"oneof=new in_progress done"`
}

type TaskStatusUpdateRequest struct {
	ID     int    `json:"id" validate:"required"`
	Status string `json:"status" validate:"required,oneof=new in_progress done"`
}
