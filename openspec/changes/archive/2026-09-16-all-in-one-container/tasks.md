## 1. Backend: раздача статики через embed FS

- [x] 1.1 Создать пакет `back/internal/application/cloudbackend/web` с `//go:embed all:dist` файловой системой и файлом-заглушкой `web/dist/.gitkeep`, и убедиться, что `go build ./cmd/cloudbackend` компилируется локально без собранного фронтенда.
- [x] 1.2 Реализовать в fiber маршрутизацию статики: `/api` и `/metrics` исключаются из обработчика, реальные файлы отдаются из embed FS, для остальных не-API путей возвращается `index.html` (SPA-fallback), и проверить на тестах в `back/internal/application/cloudbackend`.
- [x] 1.3 Удалить редирект `app.Get("/")` → `/api` (корень отдаёт `index.html`) и проверить, что существующие тесты маршрутов (`/api/health`, API-эндпоинты) остаются зелёными.
- [x] 1.4 Добавить табличные тесты на исключения `/api` и `/metrics`, отдачу статики и SPA-fallback с `index.html` для неизвестных не-API путей, и убедиться, что `make check` в `back` проходит.

## 2. Frontend: same-origin API

- [x] 2.1 Удалить dead-code ветку «localhost, port - 1» из `getBaseUrl()` в `front/src/services/Api.js`, оставив same-origin `/api`, и убедиться, что unit-тесты (если есть) и lint проходят в `front`.
- [x] 2.2 Убедиться, что `bun run build` собирает `front/dist`, который используется как источник для embed FS.

## 3. Единый Docker-образ

- [x] 3.1 Переписать `build/Dockerfile` (multi-stage: bun-сборка фронтенда → копирование `dist` в `back/.../web/dist` → сборка Go-бинарника → финальный alpine-слой) и удалить `build/Dockerfile_front`, и убедиться, что `docker build -f build/Dockerfile .` собирает образ успешно.
- [x] 3.2 Проверить, что в финальном слое присутствуют бинарник, goose, `migrations/`, не-root пользователь, каталог `/app/data` и HEALTHCHECK на `http://localhost:3000/api/health`.

## 4. Деплой: один сервис и один порт

- [x] 4.1 Сократить `deploy/compose.yaml` до одного сервиса `capuchin` с единственным маппингом `3001:3000`, сохранив SQLite-окружение, `./data:/app/data`, метки и внешнюю сеть, и убедиться, что `docker compose config` валиден и не содержит порта 3002.
- [x] 4.2 Обновить `deploy/nginx_capuchin.conf` на единственный `location /` → `http://127.0.0.1:3001`, сохранив TLS-терминацию, security-заголовки и gzip.

## 5. CI/CD и документация

- [x] 5.1 Переписать `.github/workflows/build.yml` на сборку и пуш одного образа `ghcr.io/capuchinapp/capuchin:${{github.ref_name}}` из корня репозитория через `-f build/Dockerfile --build-arg APP_VERSION=...`, убрав шаги сборки/пуша фронтенда.
- [x] 5.2 Обновить README: удалить описание blue-green deployment, описать single-container развёртывание и обновление через `docker compose up -d --force-recreate`.
