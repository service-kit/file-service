# build stage
FROM registry.ainirobot.com/arch_ci/golang:build as builder
WORKDIR /go/src/github.com/service-kit/file-service
ADD . .
RUN dep ensure -vendor-only -v
RUN go build -tags netgo -o file-service

# final stage
FROM alpine:latest
WORKDIR /file-service
RUN mkdir /file-service/conf
COPY --from=builder /go/src/github.com/service-kit/file-service /file-service/
ADD ./conf/file_service_conf.ini /file_service/conf/
RUN apk add tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo "Asia/Shanghai" >  /etc/timezone
EXPOSE 80
EXPOSE 8080
ENTRYPOINT ["./file-service"]
