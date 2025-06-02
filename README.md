# lesta-tf-idf-

## Запуск 
### Версия 0
```bash
sudo docker build -t text-analyzer .
sudo docker run -d -p 8080:8080 --name analyzer text-analyzer
```

### Версия 1
```bash
mkdir -p nginx/conf.d certbot/www certbot/conf

docker run -it --rm \
  -v $(pwd)/certbot/www:/var/www/certbot \
  -v $(pwd)/certbot/conf:/etc/letsencrypt \
  certbot/certbot certonly \
  --webroot -w /var/www/certbot \
  --email your@email.com \
  --agree-tos \
  --no-eff-email \
  -d yourdomain.com \
  -d www.yourdomain.com




openssl dhparam -out certbot/conf/ssl-dhparams.pem 2048

docker-compose up -d
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