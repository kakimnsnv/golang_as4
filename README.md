# Compose up
```
    docker-compose up -d
```

## Migrate 
```
    make migrate-up
```

### to run
```
    make run
```

### how i generated tls certificates
```
openssl req -x509 -newkey rsa:4096 -keyout server.key -out server.crt -days 365 -nodes
```