mock:
	@echo "mock repositories"
	@mockgen --source=./internal/repository/repository.go -destination=./mock/repository/repository.go --package=mock_repository

test:
	@mkdir -p coverage
	@go test -v -cover -coverpkg=./internal/service/... -coverprofile=./coverage/cover.out ./internal/service/...
	go tool cover -func=./coverage/cover.out
	@go tool cover -html=./coverage/cover.out  -o ./coverage/cover.html

.PHONY: mock test coverage