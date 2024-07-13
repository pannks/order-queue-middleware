FROM golang:1.22 as build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY main.go .

RUN go build -o main main.go

FROM public.ecr.aws/lambda/provided:al2023

COPY --from=build /app/main ./main

COPY .env .env

ENTRYPOINT ["./main"]