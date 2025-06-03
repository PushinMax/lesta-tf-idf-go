# lesta-tf-idf-

## Запуск 
### Версия 0
```bash
sudo docker build -t text-analyzer .
sudo docker run -d -p 8080:8080 --name analyzer text-analyzer
```

### Версия 1
```bash
docker-compose up -d --build
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