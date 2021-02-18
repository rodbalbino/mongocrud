# MONGODB CRUD WITH GOLANG

## Run with Docker Compose 
```
docker-compose  up -d --build
```

## OR running manualy 

### Run MongoDB
```
docker run --name mongodb -p 27017:27017 -d mongo
```

### Run Application 

```
go run *go .
```

### Load Tests 
```
cd tests
./load-test.sh
```

### Test the API

[![Run in Postman](https://run.pstmn.io/button.svg)](https://app.getpostman.com/run-collection/4d105b6060f05d5ce0b4)