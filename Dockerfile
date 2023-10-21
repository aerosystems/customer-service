FROM alpine:latest
RUN mkdir /app
RUN mkdir /app/logs

COPY ./subs-service/subs-service.bin /app

# Run the server executable
CMD [ "/app/subs-service.bin" ]