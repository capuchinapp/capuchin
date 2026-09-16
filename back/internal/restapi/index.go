package restapi

// index представляет данные для GET /.
type index struct {
	Name                   string  `json:"name"`
	IsAuth                 bool    `json:"isAuth"`
	AppVersion             string  `json:"appVersion"`
	RunningTimelogDatetime *string `json:"runningTimelogDatetime"` // format: 2006-01-02 15:04:05
}
