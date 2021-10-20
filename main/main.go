// main.go
package main

import (
	"hello/httpserver"
	"log"
	"os"
)


const dbFileName = "game.db.json"
func main() {
	db, err := os.OpenFile(dbFileName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil{
		log.Fatalf("problem opening %s %v", dbFileName, err)
	}

	store, err := httpserver.NewFileSystemPlayerStore(db)

	if err != nil {
		log.Fatalf("problem creating file system player store, %v ", err)
	}

	server := httpserver.NewPlayerServer(store)

	server.ListenAndServe()
	//if err := http.ListenAndServe(":5000", server); err != nil {
	//	log.Fatalf("cound not listen on port 5000 %v", err)
	//}





}