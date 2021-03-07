FROM golang:1.16 AS go-build

COPY ./ /app

WORKDIR /app
RUN go get
RUN go build -o main -v ./...

FROM golang:1.16 AS go-runtime

COPY --from=go-build /app/main /app/

CMD [ "/app/main" ]
