GO_FILES = $(shell find . -type f -name '*.go' | grep -v cmd/gobot/wire_gen.go)

.PHONY: test
test: $(GO_FILES)
	go test ./...

gobot: cmd/gobot/wire_gen.go $(GO_FILES)
	go build -o gobot cmd/gobot/wire_gen.go

cmd/gobot/wire_gen.go: $(GO_FILES)
	$(GOPATH)/bin/wire cmd/gobot/main.go

gobot.zip: gobot
	zip gobot.zip gobot

.PHONY: lambda
lambda: gobot.zip
	aws lambda update-function-code --function-name gobot --zip-file fileb://gobot.zip

.PHONY: clean
clean:
	git clean -Xdff

.PHONY: would-clean
would-clean:
	git clean -Xdn
