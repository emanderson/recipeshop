package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var port = flag.String("port", "8081", "port to serve on")

func RecipeShopServer(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, message())
}

func main() {
	flag.Parse()

	http.HandleFunc("/", RecipeShopServer)

	err := http.ListenAndServe(":" + *port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
