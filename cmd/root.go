package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/polarhive/wombat/crawler"
	"github.com/polarhive/wombat/db"
	"github.com/spf13/cobra"
)

var (
	seedURL string
	depth   int
	dbPath  string
)

var rootCmd = &cobra.Command{
	Use:   "wombat",
	Short: "Wombat is a Wikipedia crawler built in Go to visualize relationships between Wikipedia articles.",
	Long: `Wombat crawls Wikipedia articles, extracts links, and stores relationships between them
in a database. You can specify a seed URL, crawling depth, and the database path.`,
	Run: func(cmd *cobra.Command, args []string) {
		if seedURL == "" {
			fmt.Println("Seed URL is required!")
			os.Exit(1)
		}

		log.Println("Initializing database...")
		db.InitializeSQLiteDB(dbPath)
		defer db.CloseDB()

		log.Println("Starting to fetch and store links from seed URL:", seedURL)
		crawler.FetchAndStoreLinks(seedURL, depth, dbPath)
	},
}

func init() {
	rootCmd.Flags().StringVar(&seedURL, "seed", "", "Seed URL for crawling (required)")
	rootCmd.Flags().IntVar(&depth, "depth", 1, "Depth of the crawl (recursively)")
	rootCmd.Flags().StringVar(&dbPath, "db", "graph.db", "Path to the SQLite database")

	// Set required flag for seed
	rootCmd.MarkFlagRequired("seed")

	// Customize the help message
	rootCmd.SetHelpTemplate(`Usage: wombat [flags]

Wombat is a Wikipedia crawler that helps visualize relationships between links in Wikipedia articles.

Flags:
{{.Flags.FlagUsages | trimTrailingWhitespaces}}

Examples:
wombat --seed "https://en.wikipedia.org/wiki/Wombat" --depth 2 --db /path/to/database.db
wombat --help

For more information, visit https://github.com/polarhive/wombat
`)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
