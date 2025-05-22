mock:
	@echo "mock repositories"
	@mkdir -p mocks/repository
	@mockgen xyz/internal/repository UserRepository > mocks/repository/user_repository.go
	@mockgen xyz/internal/repository TenorLimitsRepository > mocks/repository/tenor_limits_repository.go

test:
	@mkdir -p coverage
	@go test -v -cover -coverpkg=./internal/service/... -coverprofile=./coverage/cover.out ./internal/service/...
	go tool cover -func=./coverage/cover.out
	@go tool cover -html=./coverage/cover.out  -o ./coverage/cover.html

.PHONY: mock test coverage