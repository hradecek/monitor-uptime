GO=go
AWS_BUILD_LAMBDA_ZIP=$(USERPROFILE)/Go/bin/build-lambda-zip.exe
BUILD_DIR=build
OUT_BIN=$(BUILD_DIR)/main
OUT_ZIP=$(BUILD_DIR)/main.zip
GOMAIN=lambda/main.go
SONAR_SCANNER=sonar-scanner.bat
SONAR_PROJECT_KEY=monitor-uptime
GO_TEST_JSON_REPORT=$(BUILD_DIR)/test_report.json
GO_TEST_COVERAGE_OUT=$(BUILD_DIR)/coverage.out

build: $(GOMAIN)
	@echo "> Building application into '$(OUT_BIN)'"
	GOOS=linux $(GO) build -o $(OUT_BIN) $(GOMAIN)

zip: build
	@echo "> Creating ZIP file into '$(OUT_BIN)'"
	$(AWS_BUILD_LAMBDA_ZIP) --output $(OUT_ZIP) $(OUT_BIN)

test:
	mkdir -p build
	$(GO) test -json -coverprofile build/coverage.out ./... | tee build/test_report.json

sonar: test
	@echo "> Running SonarQube analysis"
	$(SONAR_SCANNER) -Dsonar.projectKey=$(SONAR_PROJECT_KEY) \
					 -Dsonar.projetName=$(SONAR_PROJECT_KEY) \
					 -Dsonar.test.inclusions=**/*_test.go \
					 -Dsonar.sources=. \
					 -Dsonar.go.tests.reportPaths=$(GO_TEST_JSON_REPORT) \
					 -Dsonar.go.coverage.reportPaths=$(GO_TEST_COVERAGE_OUT)

clean:
	@-rm -rf $(BUILD_DIR)
