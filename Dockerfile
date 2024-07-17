FROM alpine:3.20.0
RUN mkdir /app

COPY ./customer-service.bin /app

# Run the server executable
CMD [ "/app/customer-service.bin" ]