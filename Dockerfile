FROM scratch 
WORKDIR /app
COPY ./orderservice /app

EXPOSE 8080
ENTRYPOINT ["./orderservice"]
