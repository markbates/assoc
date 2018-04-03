default:
	dep ensure -v
	go build -v -o tsoda ./vendor/github.com/gobuffalo/pop/soda
	go build -v .
	./tsoda drop -a -d
	./tsoda create -a -d
	./tsoda migrate -d
	go test -v -timeout 10s ./...
