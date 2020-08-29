GO=go
ifdef OS
	AWS_BUILD_LAMBDA_ZIP=$(USERPROFILE)/Go/bin/build-lambda-zip.exe
else
	AWS_BUILD_LAMBDA_ZIP=$(HOME)/go/bin/build-lambda-zip
endif
BUILD_DIR=build
OUT_BIN=$(BUILD_DIR)/uptime-monitor
OUT_ZIP=$(BUILD_DIR)/uptime-monitor.zip
GOMAIN=lambda/main.go
ifdef OS
	SONAR_SCANNER=sonar-scanner.bat
else
	SONAR_SCANNER=/opt/sonar-scanner/bin/sonar-scanner
endif
SONAR_PROJECT_KEY=monitor-uptime
GO_TEST_JSON_REPORT=$(BUILD_DIR)/test_report.json
GO_TEST_COVERAGE_OUT=$(BUILD_DIR)/coverage.out

zip: build
	@echo "> Creating ZIP file into '$(OUT_BIN)'"
	$(AWS_BUILD_LAMBDA_ZIP) --output $(OUT_ZIP) $(OUT_BIN)

build: clean $(GOMAIN)
	@echo "> Building application into '$(OUT_BIN)'"
	GOOS=linux $(GO) build -o $(OUT_BIN) $(GOMAIN)

test:
	mkdir -p build
	$(GO) test -json -coverprofile build/coverage.out ./... | tee build/test_report.json

sonar: test
	@echo "> Running SonarQube analysis"
	$(SONAR_SCANNER) -Dsonar.projectKey=$(SONAR_PROJECT_KEY) \
					 -Dsonar.projectName=$(SONAR_PROJECT_KEY) \
					 -Dsonar.test.inclusions=**/*_test.go \
					 -Dsonar.sources=. \
					 -Dsonar.go.tests.reportPaths=$(GO_TEST_JSON_REPORT) \
					 -Dsonar.go.coverage.reportPaths=$(GO_TEST_COVERAGE_OUT)

clean:
	@echo "> Cleaning build directory '$(OUT_BIN)'"
	@-rm -rf $(BUILD_DIR)
