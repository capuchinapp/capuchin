# Интеграционные тесты

## !!!ВАЖНО
> Используйте уникальное значение константы userID в каждом тесте
>
> Обязательно очищайте данные после выполнения теста

## Настройка VSCode `.vscode/settings.json`:
```json
{
    "go.buildFlags": ["-tags=integration sqlite"],
    "go.testTags": "integration sqlite"
}
```
