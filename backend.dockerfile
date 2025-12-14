FROM golang:1.24rc1-alpine

WORKDIR /app

RUN apk add --no-cache make

COPY . ./

RUN make get-dependencies

RUN CGO_ENABLED=0 GOOS=linux make build

EXPOSE 8080

CMD [ "./bin/musicstreaming", "-loglevel=info" ]