APP=http-server
APP_NAME=github.com/kevindoveton/http-server
PROJECT_DIR=$(GOPATH)/src/$(APP_NAME)

build:
	go install $(APP_NAME)/cmd/$(APP)

run:
	$(APP)