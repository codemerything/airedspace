package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)


func main()

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "Hello, Welcome to AiredSpace!\n")
	})

	log.Printf("Server is running on port :8080. Go to http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	log.Fatal(err)
}
