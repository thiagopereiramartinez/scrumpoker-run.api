FROM golang:1.16 AS go-build

COPY ./ /app

WORKDIR /app
RUN go get
RUN go build -o main github.com/thiagopereiramartinez/scrumpoker-run.api

FROM golang:1.16 AS go-runtime

COPY --from=go-build /app/main /app/

CMD [ "/app/main" ]
