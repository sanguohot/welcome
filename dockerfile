FROM sanguohot/cgo:v1.12.4
WORKDIR /opt/medichain
COPY . .
#禁用cgo
#RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main -v
#使用cgo并默认使用动态链接
#RUN go build -o main -v main.go
#使用cgo并使用静态链接
#RUN go build --ldflags "-linkmode external -extldflags -static" -a -o main -v
#本地模块、cgo、全静态链接
RUN GOTEMP=$(mktemp -d) && \
    GOPATH=$GOTEMP GOPROXY=file://$PWD/modvendor go build --ldflags "-linkmode external -extldflags -static" -a -o main -v && \
    rm -rf $GOTEMP

FROM busybox:1.28
WORKDIR /root/
ENV WELCOME_PATH=/root
ENV WELCOME_TYPE=production
COPY ca-certificates.crt /etc/ssl/certs/
COPY etc/config.json ./etc/
COPY sql/tbl_file_add_event_log.sql ./sql/
COPY --from=0 /opt/medichain/main .
EXPOSE 8443/tcp
ENTRYPOINT ["./main"]
