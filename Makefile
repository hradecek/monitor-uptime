GO=go
AWS_ZIP=$(USERPROFILE)\Go\bin\build-lambda-zip.exe
GOFILES=main.go dynamodb.go uptime.go sns.go
OUT_DIR=build
OUT_BIN=$(OUT_DIR)/main
OUT_ZIP=$(OUT_DIR)/main.zip

zip: build
	@echo "> Creating ZIP file into '$(OUT_BIN)'"
	$(AWS_ZIP) --output $(OUT_ZIP) $(OUT_BIN)

build: $(GOFILES)
	@echo "> Building application into '$(OUT_BIN)'"
	GOOS=linux $(GO) build -o $(OUT_BIN) $(GOFILES)

clean:
	@-rm -rf $(OUT_DIR)
