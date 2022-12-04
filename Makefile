BIN_FILE=main.out

build:
	go build -o ${BIN_FILE} main.go

clean:
	go fmt "github.com/hryoma/lc4go"
	go clean

