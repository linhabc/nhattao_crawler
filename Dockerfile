FROM golang:latest

LABEL maintainer="linhnln"

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build main.go data_type.go util.go db.go

RUN echo "Asia/Ho_Chi_Minh" > /etc/timezone
RUN dpkg-reconfigure -f noninteractive tzdata


CMD [ "./main" ]