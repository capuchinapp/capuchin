[![audit](https://github.com/capuchinapp/capuchin/actions/workflows/audit.yml/badge.svg?branch=master)](https://github.com/capuchinapp/capuchin/actions/workflows/audit.yml)

# Capuchin

## Запуск

1. Создайте каталог проекта
    ```bash
    mkdir /opt/capuchin
    cd /opt/capuchin
    ```
2. Скачайте compose-файл
    ```bash
    curl -o compose.yaml https://raw.githubusercontent.com/capuchinapp/capuchin/refs/heads/master/deploy/compose.yaml
    ```
3. Запустите приложение
    > Версия: latest или X.Y или X.Y.Z
    ```bash
    export APP_VERSION=latest && docker compose -f ./compose.yaml up -d
    ```
4. Мониторинг доступности (API и веб-интерфейс обслуживаются одним контейнером на одном порту)
    - `GET https://domain.tld/api/health` — проверка API
    - `HEAD https://domain.tld` — проверка веб-интерфейса

## Обновление

Приложение разворачивается в одном контейнере, поэтому обновление выполняется пересозданием единственного контейнера с новой версией образа:

1. Перейдите в каталог проекта
    ```bash
    cd /opt/capuchin
    ```
2. Запустите новую версию приложения, пересоздав контейнер
    > Версия: latest или X.Y или X.Y.Z
    ```bash
    export APP_VERSION=latest
    docker compose -f ./compose.yaml pull
    docker compose -f ./compose.yaml up -d --force-recreate
    ```
3. Проверьте работоспособность новой версии
    ```bash
    docker ps
    docker logs capuchin
    curl -I http://localhost:3000
    ```

## Разработка

### Релиз

1. Запустите команду `make release`

### Порты по умолчанию

- http://127.0.0.1:5172 - backend
- http://127.0.0.1:5173 - frontend

### Подготовка

1. Backend
    1. Скопируйте файл `back/.envrc.example` с именем `back/.envrc`
    2. Откройте отдельную консоль и перейдите в каталог `back`
    3. `make init`
    4. `make run`
2. Frontend
    1. Скопируйте файл `front/.env.development.example` с именем `front/.env.development`
    2. Откройте отдельную консоль и перейдите в каталог `front`
    3. `make init`
    4. `make run`

### Обновлении версий (golang, golangci-lint и alpine)

- `back/go.mod`
    ```bash
    go 1.26.0
    ```
- `back/.tool-versions`: с обязательным указанием патча, требование asdf-manager
    ```bash
    golangci-lint 2.11.4
    ```
- `build/Dockerfile`: без указания патча, чтобы сборка была на последнем патче
    ```Dockerfile
    FROM golang:1.26-alpine3.22 AS builder
    ...
    FROM alpine:3.22
    ```
- `.github/workflows/audit.yml`
    ```bash
    - name: Set up Go
      uses: actions/setup-go@v5
      with:
        go-version: '1.26'
    ...
    - name: Run linter check
      uses: golangci/golangci-lint-action@v7
      with:
        version: v2.11.4
    ```

## БД

Приложение хранит данные в одном файле SQLite. Путь задаётся переменной `SQLITE_DB_PATH` (`./back/.envrc.example`, `deploy/compose.yaml`).

Перед первым запуском накатываются миграции (`goose up`, `GOOSE_DRIVER=sqlite3`); файл БД создаётся автоматически.

### Бэкап БД

```bash
export CAPUCHIN_DT=$(date +%Y-%m-%d_%H-%M-%S)
docker exec capuchin sqlite3 /app/data/capuchin.db ".backup /app/data/backup-${CAPUCHIN_DT}.db"
docker cp capuchin:/app/data/backup-${CAPUCHIN_DT}.db ./${CAPUCHIN_DT}.db
```

### Восстановление БД из бэкапа

```bash
export CAPUCHIN_DT=2024-11-27_10-08-49
docker cp ./${CAPUCHIN_DT}.db capuchin:/app/data/restore.db
docker exec capuchin sqlite3 /app/data/restore.db "VACUUM INTO '/app/data/capuchin.db'"
```

При остановленном писателе достаточно заменить файл БД на резервную копию и удалить файлы `-wal`/`-shm`.

## Changelog сторонних библиотек

1. Back

- [Fiber](https://github.com/gofiber/fiber/releases)

2. Front

- [Air Datepicker](https://github.com/t1m0n/air-datepicker/blob/v3/CHANGELOG.md)
- [Axios](https://github.com/axios/axios/releases)
- [Bootstrap](https://github.com/twbs/bootstrap/releases)
- [Chart.js](https://github.com/chartjs/Chart.js/releases)
- [FontAwesome](https://fontawesome.com/changelog)
- [Globals](https://github.com/sindresorhus/globals/releases)
- [Svelte](https://svelte-changelog.vercel.app/)
- [sveltekit-i18n](https://github.com/sveltekit-i18n/lib/releases)
- [Tom Select](https://github.com/orchidjs/tom-select/releases)
- [Vite](https://github.com/vitejs/vite/blob/main/packages/vite/CHANGELOG.md)
- [Vitest](https://github.com/vitest-dev/vitest/releases)
- [Yup](https://github.com/jquense/yup/blob/master/CHANGELOG.md)
