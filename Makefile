export SHELL:=/bin/bash
export SHELLOPTS:=$(if $(SHELLOPTS),$(SHELLOPTS):)pipefail:errexit
GO_FILES = $(shell find . -type f -name '*.go' | grep -v wire/wire_gen.go)

# "make" will run the unit tests.
# "make lambda" will upload gobot to AWS lambda.

.PHONY: all
all: gobot history_grabber

gobot: .tested $(GO_FILES) internal/app/gobot/wire/wire_gen.go
	go build -o gobot cmd/gobot/main.go

history_grabber: $(GO_FILES) internal/app/history_grabber/wire/wire_gen.go
	go build -o history_grabber cmd/history_grabber/main.go

.PHONY: test
test: .tested

# In order to tell whether or not the code is tested or if the tests need to be re-run,
# make needs a file timestamp.  So we create a file just to store the last tested timestamp.
.tested: $(GO_FILES) internal/app/gobot/wire/wire_gen.go internal/app/history_grabber/wire/wire_gen.go
	go test ./...
	touch .tested

internal/app/gobot/wire/wire_gen.go: $(GO_FILES) internal/app/gobot/wire/wire.go
	$(GOPATH)/bin/wire internal/app/gobot/wire/wire.go

internal/app/history_grabber/wire/wire_gen.go: $(GO_FILES) internal/app/history_grabber/wire/wire.go
	$(GOPATH)/bin/wire internal/app/history_grabber/wire/wire.go

gobot.zip: .tested gobot
	zip gobot.zip gobot

.PHONY: lambda
lambda: .tested gobot.zip
	aws lambda update-function-code --function-name gobot --zip-file fileb://gobot.zip
	aws lambda wait function-updated --function-name gobot

.PHONY: clean
clean:
	git clean -Xdff

.PHONY: would-clean
would-clean:
	git clean -Xdn
