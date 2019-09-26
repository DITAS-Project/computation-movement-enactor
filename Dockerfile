FROM golang:1.13-alpine as builder

RUN mkdir $GOPATH/src/computation-movement-enactor
COPY . $GOPATH/src/computation-movement-enactor/
WORKDIR $GOPATH/src/computation-movement-enactor

ENV GO111MODULE=on

RUN CGO_ENABLED=0 GOOS=linux go build -a -o /usr/bin/computation-movement-enactor

FROM alpine:3.10

COPY --from=builder /usr/bin/computation-movement-enactor /usr/bin/computation-movement-enactor

#Trying to run the app when the container starts
RUN ["chmod", "+x", "/usr/bin/computation-movement-enactor"]
CMD ["computation-movement-enactor"]

