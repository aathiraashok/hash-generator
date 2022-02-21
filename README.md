# hash-generator
build a tool which makes http requests and prints the address of the request along with the
MD5 hash of the response.
### Build App
```
go build -o myhttp main.go
```

### Run App
```
./myhttp adjust.com google.com facebook.com 
```

### Run Tests
```
go test ./...
```