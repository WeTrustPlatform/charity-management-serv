FROM golang:1.11.2-alpine as build-binary-files

RUN apk add --update make

ENV APP_DIR /go/src/github.com/WeTrustPlatform/charity-management-serv

ARG GIT_COMMIT=unknown

WORKDIR $APP_DIR
COPY . .
RUN GIT_COMMIT=$GIT_COMMIT make build

FROM alpine:latest
RUN apk add --no-cache ca-certificates
COPY --from=build-binary-files /go/src/github.com/WeTrustPlatform/charity-management-serv/bin/ /usr/local/bin/

CMD ["cms-server"]
