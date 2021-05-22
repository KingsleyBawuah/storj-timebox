FROM golang:1.16-alpine AS GO_BUILD
COPY . /server
WORKDIR /server
RUN go build -o /go/bin/server


FROM alpine:3.10
WORKDIR app
COPY --from=GO_BUILD /go/bin/server/ ./
EXPOSE 8080
CMD ./server