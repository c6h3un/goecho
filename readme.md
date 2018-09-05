# GOECHO
a simple go echo server.

# run
```
go install
./goecho -p 8080
# build docker image
export TAG="0.4.4"
docker build -t c6h3un/goecho:$TAG .
docker push c6h3un/goecho:$TAG
```
