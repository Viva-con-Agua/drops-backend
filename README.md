# drops-backend
drops implementation with go and echo. Can handle Users.

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
