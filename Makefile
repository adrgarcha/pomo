.PHONY: build install clean test cross-compile deps

# Download dependencies
deps:
	go mod download
	go mod tidy

# Build for current platform
build: deps
	go build -o pomo main.go

# Build for all platforms
cross-compile: deps
	mkdir -p builds
	# macOS
	GOOS=darwin GOARCH=amd64 go build -o builds/pomo-darwin-amd64 main.go
	GOOS=darwin GOARCH=arm64 go build -o builds/pomo-darwin-arm64 main.go
	# Linux
	GOOS=linux GOARCH=amd64 go build -o builds/pomo-linux-amd64 main.go
	GOOS=linux GOARCH=arm64 go build -o builds/pomo-linux-arm64 main.go
	# Windows
	GOOS=windows GOARCH=amd64 go build -o builds/pomo-windows-amd64.exe main.go

# Install to /usr/local/bin (macOS/Linux)
install: build
	sudo mv pomo /usr/local/bin/

# Install to user's local bin (doesn't require sudo)
install-user: build
	mkdir -p ~/bin
	mv pomo ~/bin/
	@echo "Make sure ~/bin is in your PATH"

# Clean build artifacts
clean:
	rm -f pomo pomo.exe
	rm -rf builds/

# Run tests (when you add them)
test:
	go test -v ./...

# Run the program
run:
	go run main.go
