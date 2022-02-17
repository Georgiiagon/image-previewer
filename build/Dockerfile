FROM golang:1.16.2 as build

ENV BIN_FILE /opt/image-previewer/image-previewer
ENV CODE_DIR /go/src/

WORKDIR ${CODE_DIR}

# Кэшируем слои с модулями
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ${CODE_DIR}

# Собираем статический бинарник Go (без зависимостей на Си API),
# иначе он не будет работать в alpine образе.
ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o ${BIN_FILE} cmd/*

# На выходе тонкий образ
FROM alpine:3.9

LABEL ORGANIZATION="OTUS Online Education"
LABEL SERVICE="image-previewer"
LABEL MAINTAINERS="georgiizvonov@gmail.com"

ENV BIN_FILE "/opt/image-previewer/image-previewer"
COPY --from=build ${BIN_FILE} ${BIN_FILE}
COPY .env ${CODE_DIR}

CMD ${BIN_FILE}
