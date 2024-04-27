package main

import (
	"encoding/json"
	"fmt"
	"os"
	"rcsp/internal/model"

	"github.com/go-playground/validator/v10"
	"github.com/nats-io/stan.go"
	"github.com/sirupsen/logrus"
)

func main() {
	sc, err := stan.Connect("test-cluster", "publisher", stan.NatsURL("nats://localhost:4222"))
	if err != nil {
		panic(err)
	}
	defer sc.Close()
	path := ""
	fmt.Printf("Enter the path to file: ")
	fmt.Scan(&path)
	b, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var order model.Order
	err = json.Unmarshal([]byte(b), &order)
	if err != nil {
		logrus.Info("Failed unmarshal json\n", err)
	}
	validate := validator.New()
	err = validate.Struct(order)
	if err != nil {
		logrus.Info("Invalid json\n", err)
	} else {
		sc.Publish("addNewOrder", b)
	}
}
