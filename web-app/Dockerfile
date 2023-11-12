FROM golang:1.19-buster AS build
#RUN go env -w GOPROXY=direct

WORKDIR /app
COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY main.go ./
COPY ./static/ ./static

RUN go build -o /chatterdox

FROM gcr.io/distroless/base-debian10
WORKDIR /
COPY --from=build /chatterdox /chatterdox
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/chatterdox"]