package capuchin

import (
	"io/fs"
	"net/http"
	"path"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"

	"capuchin/internal/app/capuchin/web"
)

const (
	indexFile = "index.html"

	staticCacheControl = "public, max-age=0"
)

// webDist возвращает поддерево embed FS с файлами собранного фронтенда.
func webDist() fs.FS {
	sub, err := fs.Sub(web.Dist(), "dist")
	if err != nil {
		// Каталог dist гарантирован файлом-заглушкой .gitkeep.
		panic("web: open embedded dist: " + err.Error())
	}

	return sub
}

// staticHandler раздаёт статику фронтенда из встроенной файловой системы.
//
// Пути из excludedPaths исключаются из обработчика: API и метрики обслуживаются
// собственными роутами и серверами. Для неизвестных не-API путей возвращается
// index.html (SPA-fallback).
func staticHandler(distFS fs.FS, excludedPaths []string) fiber.Handler {
	fileHandler := adaptor.HTTPHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		file := strings.TrimPrefix(path.Clean(r.URL.Path), "/")
		if file == "" || file == "." {
			file = indexFile
		}

		if !serveEmbeddedFile(w, r, distFS, file) {
			if !serveEmbeddedFile(w, r, distFS, indexFile) {
				http.NotFound(w, r)
			}
		}
	}))

	return func(c *fiber.Ctx) error {
		if isExcludedPath(c.Path(), excludedPaths) {
			return c.Next()
		}

		if method := c.Method(); method != http.MethodGet && method != http.MethodHead {
			return c.Next()
		}

		return fileHandler(c)
	}
}

// isExcludedPath возвращает true, если путь совпадает с префиксом из excludedPaths
// или начинается с него и следующей косой черты.
func isExcludedPath(p string, excludedPaths []string) bool {
	for _, prefix := range excludedPaths {
		if p == prefix || strings.HasPrefix(p, prefix+"/") {
			return true
		}
	}

	return false
}

// serveEmbeddedFile отдаёт файл из FS и возвращает true в случае успеха.
func serveEmbeddedFile(w http.ResponseWriter, r *http.Request, distFS fs.FS, file string) bool {
	info, err := fs.Stat(distFS, file)
	if err != nil || info.IsDir() {
		return false
	}

	w.Header().Set(fiber.HeaderCacheControl, staticCacheControl)

	http.ServeFileFS(w, r, distFS, file)

	return true
}
