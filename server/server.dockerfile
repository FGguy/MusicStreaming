FROM golang:1.24rc1-alpine

WORKDIR /app

COPY . ./

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-musicstreaming

EXPOSE 8080

CMD [ "/docker-musicstreaming" ]