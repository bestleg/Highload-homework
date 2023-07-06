# syntax=docker/dockerfile:1

FROM golang:1.19

WORKDIR app/losyakov-homework

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o /losyakov-homework ./cmd/api

EXPOSE 4444

CMD [ "/losyakov-homework" ]
