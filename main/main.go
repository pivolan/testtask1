package main

import (
	"flag"
	"log"
	"testtask1"
)

var port = flag.String("port", "8098", "port")
var dsn = flag.String("dsn", "host=localhost port=5432 user=postgres dbname=testtask1 sslmode=disable", "postgres connect string")

func main() {
	b := testtask1.TestTask{}
	err := b.Init(*dsn, *port)
	if err != nil {
		log.Fatalln(err)
		return
	}

}
