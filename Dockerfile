FROM debian:latest

RUN apt-get update && apt-get install -y ca-certificates openssl

COPY ./assets/ /app/assets/
COPY ./templates/ /app/templates/
COPY ./better-opening-times /app/better-opening-times

EXPOSE 8072
WORKDIR /app
CMD ["/app/better-opening-times"]