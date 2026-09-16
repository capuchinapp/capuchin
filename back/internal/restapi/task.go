package restapi

// taskIndexFilter представляет фильтр для поиска задач.
type taskIndexFilter struct {
	ProjectID               string `query:"project_id" validate:"omitempty,ulid"`
	WithoutArchivedProjects bool   `query:"filter_archived_projects" validate:"omitempty,boolean"`
}

// taskCreateInput представляет входные данные для создания задач.
type taskCreateInput struct {
	ProjectID string `json:"projectId" validate:"required,ulid"`
	Name      string `json:"name" validate:"required,max=255"`
	Comment   string `json:"comment"`
}

// taskUpdateInput представляет входные данные для обновления задач.
type taskUpdateInput struct {
	ProjectID string `json:"projectId" validate:"required,ulid"`
	Name      string `json:"name" validate:"required,max=255"`
	Comment   string `json:"comment"`
}
