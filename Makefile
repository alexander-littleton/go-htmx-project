install:
	go install github.com/a-h/templ/cmd/templ@latest

run:
	templ generate
	go run main.go