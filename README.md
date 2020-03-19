# tractor

we like to play games

## installation

- set up your go [workspace](https://medium.com/golang-learn/go-project-layout-e5213cdcfaa2): the tractor project should be in something like
  `$GOPATH/github.com/{user}/tractor`. I use the default `$GOPATH` but you can change this as desired.
- set up client:

```
brew install yarn
cd client
yarn build
```

## usage

- `make docker-build`: builds docker image, server running at `localhost:8080`
- `make docker-cleanup`: clean up docker container and image from `docker-build`
- for local dev: `make server` to start go server at `localhost:8080`, and `make app` to build client, running at `localhost:3000`
