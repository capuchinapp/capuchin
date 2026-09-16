package restapi

// timelogIndexFilter представляет фильтр для поиска записей в журнале.
type timelogIndexFilter struct {
	DateFrom  string `query:"date_from" validate:"datetime=2006-01-02"`
	DateTo    string `query:"date_to" validate:"datetime=2006-01-02"`
	ClientID  string `query:"client_id" validate:"omitempty,ulid"`
	ProjectID string `query:"project_id" validate:"omitempty,ulid"`
}

// timelogCreateInput представляет входные данные для создания записи в журнале.
type timelogCreateInput struct {
	ProjectID    string  `json:"projectId" validate:"required,ulid"`
	Date         string  `json:"date" validate:"datetime=2006-01-02"`
	TimeStart    string  `json:"timeStart" validate:"datetime=15:04:05"`
	TimeEnd      *string `json:"timeEnd" validate:"omitempty,datetime=15:04:05"`
	BillableRate int32   `json:"billableRate" validate:"min=0"`
	Comment      string  `json:"comment"`
	TaskID       *string `json:"taskId" validate:"omitempty,ulid"`
}

// timelogUpdateInput представляет входные данные для обновления записи в журнале.
type timelogUpdateInput struct {
	ProjectID    string  `json:"projectId" validate:"required,ulid"`
	Date         string  `json:"date" validate:"datetime=2006-01-02"`
	TimeStart    string  `json:"timeStart" validate:"datetime=15:04:05"`
	TimeEnd      string  `json:"timeEnd" validate:"datetime=15:04:05"`
	BillableRate int32   `json:"billableRate" validate:"min=0"`
	Comment      string  `json:"comment"`
	TaskID       *string `json:"taskId" validate:"omitempty,ulid"`
}

// timelogStopInput представляет входные данные для остановки записи в журнале.
type timelogStopInput struct {
	Date    string `json:"date" validate:"datetime=2006-01-02"`
	TimeEnd string `json:"timeEnd" validate:"datetime=15:04:05"`
}
