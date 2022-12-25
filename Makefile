.PHONY: build
build:
	go build -o bin/columnar_store

.PHONY: run
run:
	./bin/columnar_store

.PHONY: clean
clean:
	rm -f bin/columnar_store
	rm -f ~/columnar/data/meta/_tables
	rm -rf ~/columnar/data/tables/*

.PHONY: test
test: clean build
	./testit.sh

.PHONY: fake
fake:
	echo "fake"
