FROM golang:1.11.2

ENV APP_DIR /go/src/github.com/WeTrustPlatform/charity-management-serv
WORKDIR $APP_DIR
COPY . .
RUN make build
CMD make launch
