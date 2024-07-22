docker-build:
	@echo "DOCKER: Building image"
	@docker build -t pkdevel/docker-home .

docker-run:
	@echo "DOCKER: Starting container"
	@docker run --rm -d \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-p 9080:8080 \
		--name=docker-home \
		pkdevel/docker-home

generate: _gow _templ _tailwind
	@echo "TEMPL: Generating templates"
	@templ generate
	@echo "TAILWIND: Generating styles"
	@tailwindcss -c web/tailwind.config.js -i web/style/tailwind.css -o assets/style.css -m

build: generate
	@echo "GO: Building"
	@go build -v ./cmd/main.go

run: generate
	@echo "GO: Starting"
	@go run ./cmd/main.go

clean:
	@echo "GO: Cleaning"
	@go clean
	@echo "DOCKER: Cleaning"
	@docker image prune --filter label=name=docker-home --force --all
	@docker builder prune --force

go-watch: _gow
	gow run ./cmd/main.go

templ-watch: _templ
	templ generate -watch

tailwind-watch:
	tailwindcss -c web/tailwind.config.js -i web/style/tailwind.css -o assets/style.css -mw

_gow:
	@if ! command -v gow &> /dev/null; then \
		echo "GO: Installing gow"; \
		go install github.com/mitranim/gow@latest; \
	fi

_templ:
	@if ! command -v templ &> /dev/null; then \
		echo "TEMPL: Installing templ"; \
		go install github.com/a-h/templ/cmd/templ@latest; \
	fi

_tailwind:
	@if ! command -v tailwindcss &> /dev/null; then \
		echo "TAILWIND: cli not found, please install..."; \
		exit 1; \
	fi
	@# npm install -g tailwindcss; \
