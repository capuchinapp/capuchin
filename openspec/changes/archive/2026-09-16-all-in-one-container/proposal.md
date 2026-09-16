## Why

Приложение ранее состояло из трёх контейнеров (PostgreSQL, бэкенд, фронтенд) с blue-green deployment для облачной версии с неограниченным кругом пользователей. Сейчас preparations для self-hosted версии на SQLite уже выполнены: миграция с PostgreSQL завершена, лишние облачные функции удалены. Однако инфраструктура деплоя не обновлена — по-прежнему существуют два отдельных Dockerfile (Dockerfile_back, Dockerfile_front), фронтенд раздаётся через nginx в отдельном контейнере, а README описывает blue-green сценарий, который для SQLite опасен (два одновременно пишущих процесса в одну БД). Необходимо упростить всё до одного контейнера.

## What Changes

- Объединение бэкенда и фронтенда в один Docker-образ: Go-приложение само раздаёт статику фронтенда ( embed FS )
- Удаление Dockerfile_front и связанного nginx-контейнера
- Обновление Dockerfile_back → единый Dockerfile для all-in-one образа
- Билд фронтенда (bun + vite) встраивается в multi-stage Docker build как промежуточный шаг
- Упрощение compose.yaml до одного сервиса с одним портом (API + статика на одном порту)
- Удаление описания blue-green deployment из README
- Обновление deploy/nginx_capuchin.conf (это обратный прокси для выставления приложения в интернет)
- Обновление CI/CD workflow (build.yml) для сборки одного образа вместо двух

## Capabilities

### New Capabilities

- `deployment/single-container`: All-in-one Docker образ с бэкендом, фронтендом и SQLite внутри. Go-приложение раздаёт статику фронтенда и API на одном порту. Единый compose.yaml для self-hosted развёртывания.

### Modified Capabilities

(нет — ранее specs не существовало)

## Impact

- **Удаляемые файлы**: `build/Dockerfile_front`
- **Изменяемые файлы**: `build/Dockerfile_back` (rename → `build/Dockerfile`), `deploy/compose.yaml`, `back/internal/application/cloudbackend/app.go` (добавление раздачи статики), `.github/workflows/build.yml`
- **Порты**: один порт вместо двух (текущий 3000 → единый порт для API и фронтенда)
- **Обратная совместимость**:.blue-green deployment полностью удаляется; для обновления используется `docker compose up -d --force-recreate`
