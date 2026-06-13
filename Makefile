.PHONY: run build clean

run:
	go run main.go

build:
	go build -o mefetch-readme .

clean:
	rm -f mefetch-readme