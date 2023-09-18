package main

import (
	"flag"
	"log"
	"security/proxy_server"
	"security/store"
)

func main() {
	httpsFlag := flag.Bool("https", false, "Работа с защищенным соединением.")
	flag.Parse()

	storage, err := store.NewStore()
	if err != nil {
		log.Fatal(err)
	}

	server := proxy_server.NewServer(storage)

	if *httpsFlag {
		server.StartServerTLS()
	}

	server.StartServer()
}
