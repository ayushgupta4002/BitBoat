build:
	go build -o bin/main 

run: build
	./bin/main --listenaddr localhost:3000

runsubs: build
	./bin/main --listenaddr localhost:8080 --adminaddr localhost:3000