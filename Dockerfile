FROM alpine:latest
RUN mkdir /app

COPY ./customer-service.bin /app

# Run the server executable
CMD [ "/app/customer-service.bin" ]