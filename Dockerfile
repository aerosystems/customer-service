FROM alpine:latest
RUN mkdir /app
RUN mkdir /app/logs

COPY ./customer-service.bin /app

# Run the server executable
CMD [ "/app/customer-service.bin" ]