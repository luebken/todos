# TODOs

This is a simple TODO app written in Golang talking to a Postgres database.

# Getting started

See [Makefile] to get started.

## Local K8s

```sh
# create cluster
kind create cluster

# setup database
kubectl apply -f k8s/postgres.yaml
kubectl port-forward service/postgres 5432:5432
psql -h localhost -U postgres todos
CREATE TABLE todos (item TEXT PRIMARY KEY);
INSERT INTO todos (item) VALUES ('Buy groceries'), ('Finish homework'), ('Clean the house');

# start application
kubectl apply -f k8s/todos.yaml
kubectl port-forward service/todos 3000:3000

# tear down cluster
kind delete cluster
```

## Troubleshooting

```sh
# base connectivity
kubectl run netcat --rm -it --image=alpine -- sh
nc -zv postgres 5432

# psql
kubectl run psql-client --rm -it --image=postgres -- bash
psql -h postgres -U postgres 
```

# Notes
Inspired by https://blog.logrocket.com/building-simple-app-go-postgresql/