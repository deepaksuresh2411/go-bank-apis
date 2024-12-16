package main

import "log"

func main() {
	storage, err := NewPostgresStorage()
	if err != nil {
		log.Fatal(err)
	}
	if err := storage.CreateAccountTable(); err != nil {
		log.Fatal(err)
	}
	server := NewAPIServer(":8000", storage)
	server.Start()
}
