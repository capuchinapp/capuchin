package restapi

// favoriteInput представляет входные данные для создания/обновления избранного.
type favoriteInput struct {
	Name         string  `json:"name" validate:"required,max=255"`
	ProjectID    string  `json:"projectId" validate:"required,ulid"`
	TaskID       *string `json:"taskId" validate:"omitempty,ulid"`
	BillableRate int32   `json:"billableRate" validate:"min=0"`
	Comment      string  `json:"comment"`
}
