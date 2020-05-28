run:
	go run main.go

docker:
	docker build -t msawady/yf-bot-go:latest .

run-docker:
	docker run msawady/yf-bot-go

build:
	go build

clean:
	go clean
	rm -rf yf-bot-go
