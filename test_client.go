package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func main() {
	s := struct {
		Path  string
		Delay bool
	}{
		"http://ecx.images-amazon.com/images/I/51eDwv7tCtL._SX442_BO1,204,203,200_.jpg",
		false,
	}

	data, err := json.Marshal(s)
	if err != nil {
		os.Exit(1)
	}

	resp, err := http.Post("http://localhost:8080/thumbnail", "application/json", bytes.NewBuffer(data))
	if err != nil {
		fmt.Println("Get error ", err)
	} else {
		fmt.Println(resp.Status)
	}
}
