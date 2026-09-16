package restapi

// authRegisterInput представляет входные данные для регистрации.
type authRegisterInput struct {
	Email string `json:"email" validate:"required,email"`
}

// authActivateInput представляет входные данные для активации.
type authActivateInput struct {
	UserID string `json:"userId" validate:"required,ulid"`
	Code   string `json:"code" validate:"required,len=40"` // См. тип поля в файле migrations/00002_alter__users__activate_code.sql
}

// authLoginInput представляет входные данные для авторизации.
type authLoginInput struct {
	Email string `json:"email" validate:"required,email"`
}

// authApplyCodeInput представляет входные данные для применения кода.
type authApplyCodeInput struct {
	Email string `json:"email" validate:"required,email"`
	Code  string `json:"code" validate:"required,len=6"`
}
