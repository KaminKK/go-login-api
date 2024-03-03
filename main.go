package main

import (
	"fmt"
	"net/http"
)

func main() {

	db, err := ConnectDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	http.HandleFunc("/register", RegisterHandler)
	http.HandleFunc("/login", LoginHandler)

	fmt.Println("Server is running on localhost:8080")
	http.ListenAndServe(":8080", nil)
}
