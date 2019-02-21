FROM golang:1.11.2-alpine

RUN apk add --update make

ENV APP_DIR /go/src/github.com/WeTrustPlatform/charity-management-serv

ARG GIT_COMMIT=unknown

WORKDIR $APP_DIR
COPY . .
RUN GIT_COMMIT=$GIT_COMMIT make build
CMD make launch
