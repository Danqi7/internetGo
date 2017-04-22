package main

import (
	"fmt"
	"net/http"
)

func main() {
	response, err := http.Get("http://192.168.1.103:3128")
	if err != nil {
		fmt.Println(err)
        return
	}
	fmt.Println(response.Body)

}
