// Package web содержит собранный фронтенд, встроенный в бинарник через embed FS.
package web

import "embed"

//go:embed all:dist
var distFiles embed.FS

// Dist возвращает embed FS с файлами собранного фронтенда.
func Dist() embed.FS {
	return distFiles
}
