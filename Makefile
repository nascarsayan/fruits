VER ?= 4

run:
	go run main.go

docker-build:
	docker build -t fruits:$(VER) -f refs/Dockerfile.$(VER) .

docker-run:
	docker run -it -p 9999:9999 fruits:$(VER)

zip:
	git archive --format zip --output fruits.zip master
