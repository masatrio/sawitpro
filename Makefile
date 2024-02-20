.PHONY: clean all init generate generate_mocks

all: build/main

build/main: cmd/main.go generated
	@echo "Building..."
	go build -o $@ $<

clean:
	rm -rf generated

init: generate
	go mod tidy
	go mod vendor

test:
	go test -tags="!test" -short -coverprofile coverage.out -v ./...

generate: generated generate_mocks

generated: api.yml
	@echo "Generating files..."
	mkdir generated || true
	oapi-codegen --package generated -generate types,server,spec $< > generated/api.gen.go

COMMON_INTERFACE := common/interfaces.go
REPOSITORY_INTERFACE := repository/interfaces.go
SERVICE_INTERFACE := service/interfaces.go

COMMON_MOCK := mocks/common_mock.gen.go
REPOSITORY_MOCK := mocks/repository_mock.gen.go
SERVICE_MOCK := mocks/service_mock.gen.go

generate_mocks: $(COMMON_MOCK) $(REPOSITORY_MOCK) $(SERVICE_MOCK)

$(COMMON_MOCK): $(COMMON_INTERFACE)
	@echo "Generating mocks for common interfaces..."
	mockgen -source=$< -destination=$@ -package=mocks

$(REPOSITORY_MOCK): $(REPOSITORY_INTERFACE)
	@echo "Generating mocks for repository interfaces..."
	mockgen -source=$< -destination=$@ -package=mocks

$(SERVICE_MOCK): $(SERVICE_INTERFACE)
	@echo "Generating mocks for service interfaces..."
	mockgen -source=$< -destination=$@ -package=mocks