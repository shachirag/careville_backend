# Build Stage
FROM golang:1.20 AS build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -o main main.go

# Run Stage
FROM gcr.io/distroless/static-debian11
COPY --from=build /app/main /
COPY ./fkpmobileapplication-firebase-adminsdk-fnutv-be92c67c9c.json .
# COPY ./public /public

EXPOSE 8080
CMD ["/main"]
