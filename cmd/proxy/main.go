package main

import (
	"flag"
	"security/proxy_server"
)

func main() {
	httpsFlag := flag.Bool("https", false, "Работа с защищенным соединением.")
	flag.Parse()

	if *httpsFlag {
		proxy_server.StartServerTLS()
	}

	proxy_server.StartServer()
}
