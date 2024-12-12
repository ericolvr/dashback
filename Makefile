# ANSI color codes
COLOR_RESET=\033[0m
COLOR_BOLD=\033[1m
COLOR_GREEN=\033[32m
COLOR_YELLOW=\033[33m

help:
	@echo ""
	@echo "  $(COLOR_YELLOW)Available targets:$(COLOR_RESET)"
	@echo ""
	@echo "  $(COLOR_GREEN)containers$(COLOR_RESET)		- Start Containers - Rabbit - Mongodb"
	@echo "  $(COLOR_GREEN)install$(COLOR_RESET)		- Install Dependencies"
	@echo "  $(COLOR_GREEN)run$(COLOR_RESET)			- Run Development Server on Localhost"
	@echo "  $(COLOR_GREEN)build$(COLOR_RESET)			- Build and Run to Production"
	@echo ""
	@echo "  $(COLOR_YELLOW)Note:$(COLOR_RESET) Use 'make <target>' to execute a specific target."
	@echo ""

install:
	go mod tidy

containers:
	docker-compose up -d

run:
	go run cmd/main.go

build:
	go build -o bin/main cmd/main.go && ./bin/main &


.PHONY: containers, install, run, build
