package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
)

type EventBody struct {
	Message struct {
		Data []byte `json:"data"`
	} `json:"message"`
	Subscription string `json:"subscription"`
}

type Customer struct {
	Email string `json:"email"`
	UID   string `json:"uid"`
}

func main() {
	customer := Customer{
		Email: "customer@gmail.com",
		UID:   uuid.New().String(),
	}

	data, err := json.Marshal(customer)
	if err != nil {
		panic(err)
	}

	var event EventBody
	event.Message.Data = data
	event.Subscription = "create-customer"

	var jsonData []byte
	jsonData, err = json.Marshal(event)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(jsonData))
}
