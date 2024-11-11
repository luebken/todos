# Emma Demo

Show color in the terminal:
```
kubecolor get pods -o wide -n todos --force-colors
watch -n 1 --color kubecolor get pods -o wide -n todos  --force-colors
```

Label notes:
```
kubectl label --list nodes hgywdr-worker-node-azure kubectl label --list nodes musdav-worker-node-aws kubectl label --list nodes ntrfaw-worker-node-gcp

emma.ms/cloud-provider=azure emma.ms/location=germany,frankfurt emma.ms/data-center=europe-west3-b
```

ingress
```
k -n ingress exec -ti emma-ingress-controller-ingress-nginx-controller-6b6c57d4ck4rsn -- cat /etc/nginx/nginx.conf```
```


TODOS
* read password from the same secret
* script to move postgres including copying pg_dump production-db | psql test-db


