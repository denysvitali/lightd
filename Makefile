TARGET_IP?=192.168.1.166

build:
	CGO_ENABLED=0 GOARCH=arm64 go build -o bin/lightd

deploy:
	scp bin/lightd "root@$(TARGET_IP)":/opt/bin/lightd

.PHONY: build
