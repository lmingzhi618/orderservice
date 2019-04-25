FROM scratch 
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY ./orderservice /app

EXPOSE 8080
ENTRYPOINT ["./orderservice"]
