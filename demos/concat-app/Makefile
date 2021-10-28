all: concat-app

concat-app: *.go
	go build -o $@ $^

clean:
	rm -f concat-app
