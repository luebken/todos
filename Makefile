dev:
	go run cmd/todos/server.go

build:
	docker build . -f Dockerfile -t luebken/todos:latest

run:
	docker run -p 3000:3000 luebken/todos:latest

push: build
	docker push luebken/todos:latest
	kubectl delete pods -l app=todos

run-postgres:
	docker run -p 5432:5432 -e POSTGRES_PASSWORD=mysecretpassword postgres