FROM golang:1.16 as builder

LABEL maintainer="lachezarbogomilov@gmail.com"

ENV BASE_APP_DIR $GOPATH/go/src/github.com/la4ezar/restapi
ENV CGO_ENABLED 0
ENV GOOS linux
WORKDIR ${BASE_APP_DIR}

COPY go.mod ${BASE_APP_DIR}

#RUN go mod download

COPY . ${BASE_APP_DIR}

RUN go build -o /app/server cmd/server/server.go

COPY application.yaml /app
#\
#&& go build -o /app/client cmd/client/client.go


FROM alpine

EXPOSE 8080 8080

WORKDIR /app

COPY --from=builder /app /app

#CMD ["/app/server", "&", "/app/client"]
CMD ["/app/server"]
