FROM obraun/vss-protoactor-jenkins as builder
COPY . /app
WORKDIR /app
RUN go build -o treeservice/main treeservice/main.go

FROM iron/go
COPY --from=builder /app/treeservice/main /app/treeservice
EXPOSE 8091
ENTRYPOINT ["/app/treeservice"]
