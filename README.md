# drops-backend
drops implementation with go and echo. Can handle Users.


## .env

```
REPO_PORT=1323
IROBERT_URL="https://irobertstage.vivaconagua.org"

CRM_SIGNUP=false

REPO_CONFIG_PATH=./config
VOLUME_PATH=/home/dls/Workspace/volumes

ALLOW_ORIGINS=http://localhost:8080,http://172.2.60.1
COOKIE_SECURE=false
SAME_SITE=none


DB_HOST=localhost
DB_PORT=27018
DB_NAME=drops


MYSQL_DATABASE=drops
MYSQL_USER=drops
MYSQL_PASSWORD=drops 
MYSQL_ROOT_PASSWORD=yes 

NATS_HOST=localhost
NATS_PORT=4222

REDIS_HOST=localhost
REDIS_PORT=6379

DROPS_HOST=localhost
DROPS_PORT=1323
```

## use
install docker-compose

## redis
run `docker-compose up -d redis`
## database
run `docker-compose up -d drops-database && mysql -u drops -pdrops -h 172.2.200.1 drops < drops-database.sql`

## development

### 1. Install Go language 
Like here: https://itrig.de/index.php?/archives/2377-Installation-einer-aktuellen-Go-Version-auf-Ubuntu.html

### 2. Install dependecies

~~~
go get github.com/go-playground/validator
go get github.com/go-redis/redis
go get github.com/go-sql-driver/mysql
go get github.com/google/uuid
go get github.com/jinzhu/configor
go get github.com/labstack/echo
go get github.com/labstack/echo-contrib/session
go get github.com/rbcervilla/redisstore
go get golang.org/x/crypto/bcrypt
~~~

### 3. Checkout drops-backend
git clone https://github.com/Viva-con-Agua/drops-backend.git

### 4. Run server
Start server wiht `go run server.go`

### 5. update nginx
Update IP to you local IP in develop-pool branch at `routes/nginx-pool/pool.upstream` and restart nginx-pool docker with `docker restart pool-nginx`
