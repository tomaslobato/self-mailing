.PHONY: install

install:
	go build -o self-mailing ./main.go
	sudo mv self-mailing /usr/local/bin/self-mailing
