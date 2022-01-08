FROM atlassian/ssh-ubuntu

COPY id_rsa ~/.ssh

RUN go env -w GOPROXY=https://goproxy.cn,direct

RUN go build -o zssh

ENTRYPOINT ["./zssh"]