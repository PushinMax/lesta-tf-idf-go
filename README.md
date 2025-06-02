# lesta-tf-idf-

## Запуск 
```bash
docker build -t text-analyzer .
docker run -d -p 8080:8080 --name analyzer text-analyzer
```

## Для разработки
### Убить процесс
```bash
sudo kill -9 $(sudo lsof -t -i:8080)
```

### Подключиться в VM
```bash
ssh ubuntu@37.9.53.117 -i ~/.ssh/private.key
```