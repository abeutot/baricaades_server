FROM golang:latest

WORKDIR /go/src/github.com/abeutot/baricaades_server
COPY . .

RUN go test ./...
RUN CGO_ENABLED=0 go build .

FROM alpine:latest

ARG CORS_ALLOW_ORIGINS=http://localhost:3000

ENV GIN_MODE=release
ENV CORS_ALLOW_ORIGINS=$CORS_ALLOW_ORIGINS

EXPOSE 8080/tcp

WORKDIR /root/
COPY --from=0 /go/src/github.com/abeutot/baricaades_server/baricaades_server .

CMD ["./baricaades_server"]
