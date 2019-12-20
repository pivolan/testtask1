package main

import (
	"flag"
	"log"
	"testtask1"
)

var port = flag.String("port", "8098", "port")
var db = flag.String("db", "host=localhost port=5432 user=postgres dbname=testtask1 password=postgres", "postgres connect string")

func main() {
	b := testtask1.TestTask{}
	err := b.Init(*db, *port)
	if err != nil {
		log.Fatalln(err)
		return
	}
}
