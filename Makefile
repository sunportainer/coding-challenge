PLATFORM=$(shell go env GOOS)
ARCH=$(shell go env GOARCH)

ifeq ("$(PLATFORM)", "windows")
bin=app.exe
else
bin=app
endif
dist := dist

.PHONY: pre $(bin) test clean

all: test $(bin)

pre:
	mkdir -pv $(dist)
test:
	go test ./service

$(bin): pre
	GOOS="$(shell go env GOOS)" GOARCH="$(shell go env GOARCH)" CGO_ENABLED=0 go build -a --installsuffix cgo --ldflags '-s' -o dist/$(bin)

clean:
	rm -rf $(dist)/*