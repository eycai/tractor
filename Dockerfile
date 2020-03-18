# Build the Go API
FROM golang:latest AS builder
ENV GO111MODULE=on
WORKDIR /go/src
COPY src/go.mod ./
COPY src/go.sum .
RUN go mod download
COPY src ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w" -a -o /main cmd/tractor/main.go

# Build the React application
FROM node:alpine AS node_builder
COPY client ./
RUN yarn install
RUN yarn build

# Final stage build, this will be the container
# that we will deploy to production
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /main ./
COPY --from=node_builder /build ./web
RUN chmod +x ./main
EXPOSE 8080
CMD ./main