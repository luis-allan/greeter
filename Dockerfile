# build stage
FROM golang:alpine AS build-env
RUN apk --no-cache add build-base bash git
ADD . /go/src/github.com/thingful/greeter
WORKDIR /go/src/github.com/thingful/greeter
RUN make build-internal

# final stage
FROM alpine
RUN apk --no-cache add ca-certificates postgresql-client && rm -rf /var/cache/apk/*
WORKDIR /app
COPY --from=build-env /go/src/github.com/thingful/greeter/build/greeter /app/
ENTRYPOINT [ "/app/greeter" ]
CMD [ "" ]
