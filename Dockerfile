FROM scratch 
WORKDIR /app
COPY . /app

EXPOSE 8080
ENTRYPOINT ["./orderservice"]
