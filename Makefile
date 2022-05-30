build:
	echo "Compiling"
	GOOS=linux GOARCH=amd64 go build -o bin/recipes-amd64-linux cli.go
	GOOS=windows GOARCH=amd64 go build -o bin/recipes-amd64-windows cli.go
	GOOS=darwin GOARCH=amd64 go build -o bin/recipes-amd64-darwin cli.go