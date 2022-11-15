package main

import (
	"encoding/json"
	"net"
)

type data struct {
	Foo string `json:"foor"`
	Bar int    `json:"bar"`
}

func main() {
	d := data{
		Foo: "LALALLA",
		Bar: 1337,
	}

	data, err := json.Marshal(&d)
	if err != nil {
		panic(err)
	}

	listen, err := net.Listen("tcp", "0.0.0.0:9999")
	if err != nil {
		panic(err)
	}
	defer listen.Close()

	for {
		c, err := listen.Accept()
		if err != nil {
			panic(err)
		}
		c.Write(data)
		c.Close()
	}
}
