package restapi

// projectIndexFilter представляет фильтр для поиска проектов.
type projectIndexFilter struct {
	ClientID               string `query:"client_id" validate:"omitempty,ulid"`
	WithoutArchivedClients bool   `query:"filter_archived_clients" validate:"omitempty,boolean"`
}

// projectCreateInput представляет входные данные для создания проекта.
type projectCreateInput struct {
	ClientID     string `json:"clientId" validate:"required,ulid"`
	Name         string `json:"name" validate:"required,max=255"`
	BillableRate int32  `json:"billableRate" validate:"min=0"`
	Comment      string `json:"comment"`
}

// projectUpdateInput представляет входные данные для обновления проекта.
type projectUpdateInput struct {
	ClientID     string `json:"clientId" validate:"required,ulid"`
	Name         string `json:"name" validate:"required,max=255"`
	BillableRate int32  `json:"billableRate" validate:"min=0"`
	Comment      string `json:"comment"`
}
