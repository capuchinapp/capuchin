package restapi

// settingUpdateInput представляет входные данные для обновления настроек.
type settingUpdateInput struct {
	// Формат даты.
	// Пример: DD.MM.YYYY, MM/DD/YYYY
	// https://day.js.org/docs/en/display/format
	DateFormat string `json:"dateFormat" validate:"len=10"`

	// Перечисление рабочих дней недели.
	// Пример: 1,2,3,4,5
	WorkingDays string `json:"workingDays"`
}
