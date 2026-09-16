package sqlite

import "capuchin/internal/domain"

// settingModel представляет настройку.
type settingModel struct {
	UserID string `db:"user_id"`
	Key    string `db:"key"`
	Value  string `db:"value"`
}

// settingModelFromDomain преобразует настройку домена в настройку sqlite.
func settingModelFromDomain(src domain.Setting) settingModel {
	return settingModel{
		UserID: src.UserID,
		Key:    src.Key,
		Value:  src.Value,
	}
}

// settingModelToDomain преобразует настройку sqlite в настройку домена.
func settingModelToDomain(src settingModel) domain.Setting {
	return domain.Setting{
		UserID: src.UserID,
		Key:    src.Key,
		Value:  src.Value,
	}
}

// settingModelsToDomains преобразует настройки sqlite в настройки домена.
func settingModelsToDomains(src []settingModel) []domain.Setting {
	dst := make([]domain.Setting, len(src))

	for i := range src {
		dst[i] = settingModelToDomain(src[i])
	}

	return dst
}
