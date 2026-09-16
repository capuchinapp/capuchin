package restapi

// clientCreateInput представляет входные данные для создания клиента.
type clientCreateInput struct {
	Name         string `json:"name" validate:"required,max=255"`
	BillableRate int32  `json:"billableRate" validate:"min=0"`
	Comment      string `json:"comment"`
}

// clientUpdateInput представляет входные данные для обновления клиента.
type clientUpdateInput struct {
	Name         string `json:"name" validate:"required,max=255"`
	BillableRate int32  `json:"billableRate" validate:"min=0"`
	Comment      string `json:"comment"`
}
