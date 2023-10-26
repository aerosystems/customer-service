FROM alpine:latest
RUN mkdir /app
RUN mkdir /app/logs

COPY ./user-service/user-service.bin /app

# Run the server executable
CMD [ "/app/user-service.bin" ]