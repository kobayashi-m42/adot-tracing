FROM golang:1.20

WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o application .

EXPOSE 5001

CMD ["./application"]
