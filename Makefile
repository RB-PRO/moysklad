all: run

run:
	go run cmd/main/main.go

push:
	git push git@github.com:RB-PRO/moysklad.git

pull:
	git pull git@github.com:RB-PRO/moysklad.git

pushW:
	git push https://github.com/RB-PRO/moysklad.git

pullW:
	git pull https://github.com/RB-PRO/moysklad.git

build-config:
	go env GOOS GOARCH

build-linux:
	set GOARCH=amd64
	set GOOS=linux
	go build cmd/main/main.go

build-windows:
	set GOARCH=amd64
	set GOOS=windows
	go build cmd/main/main.go