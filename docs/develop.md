# Для разработки
## Убить процесс
```bash
sudo kill -9 $(sudo lsof -t -i:8080)
```

## Подключиться к VM
```bash
ssh ubuntu@37.9.53.117 -i ~/.ssh/private.key
```

## Генерация .env файла
```bash
sh scripts/generate_env.sh
```

## **Подключение к базе данных**
```bash
docker exec -it postgres_db psql -U postgres