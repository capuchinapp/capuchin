package domain

// Favorite представляет элемент избранного.
type Favorite struct {
	ID           string  `json:"id"`
	UserID       string  `json:"-"`
	Name         string  `json:"name"`
	ProjectID    string  `json:"projectId"`
	ProjectName  string  `json:"projectName"`
	ClientID     string  `json:"clientId"`
	ClientName   string  `json:"clientName"`
	TaskID       *string `json:"taskId"`
	TaskName     *string `json:"taskName"`
	BillableRate int32   `json:"billableRate"`
	Comment      string  `json:"comment"`
}
