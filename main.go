package main

import (
	"log"

	"github.com/polarhive/wombat/crawler"
	"github.com/polarhive/wombat/db"
)

const seedURL = "https://en.wikipedia.org/wiki/Palm_nut"

func main() {
	log.Println("Initializing database...")
	db.InitializeSQLiteDB("graph.db")
	defer db.CloseDB()

	log.Println("Starting to fetch and store links from seed URL:", seedURL)
	crawler.FetchAndStoreLinks(seedURL)
}
