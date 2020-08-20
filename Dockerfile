FROM golang:latest

LABEL maintainer="linhnln"

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build main.go data_type.go util.go db.go

ENV TZ=Asia/Ho_Chi_Minh
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

CMD [ "./main" ]