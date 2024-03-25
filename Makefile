install:
	go install github.com/a-h/templ/cmd/templ@latest

run:
	templ generate --watch --proxy="http://localhost:3000" --cmd="go run ."