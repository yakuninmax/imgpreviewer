FROM golang:1.22.6-alpine3.20 as build

ENV BIN_FILE /opt/imgpreviewer/imgpreviewer
ENV CODE_DIR /go/src/

WORKDIR ${CODE_DIR}

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ${CODE_DIR}

RUN CGO_ENABLED=0 go build -o ${BIN_FILE} cmd/imgpreviewer/*

FROM alpine:3.20.3

LABEL ORGANIZATION="OTUS Online Education"
LABEL SERVICE="image previewer"
LABEL MAINTAINERS="yakunin.max@gmail.com"

ENV BIN_FILE "/opt/imgpreviewer/imgpreviewer"
COPY --from=build ${BIN_FILE} ${BIN_FILE}

CMD ${BIN_FILE}
