package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/polarhive/wombat/crawler"
	"github.com/polarhive/wombat/db"
	"github.com/polarhive/wombat/export"
	"github.com/spf13/cobra"
)

var (
	seedURL  string
	depth    int
	dbPath   string
	exportTo string
)

var rootCmd = &cobra.Command{
	Use:   "wombat",
	Short: "Wombat is a Wikipedia crawler built in Go to visualize relationships between Wikipedia articles.",
	Long: `Wombat crawls Wikipedia articles, extracts links, and stores relationships between them
in a database. You can specify a seed URL, crawling depth, and the database path.`,

	Run: func(cmd *cobra.Command, args []string) {
		if exportTo == "md" {
			if dbPath == "" {
				fmt.Println("Database path is required for export!")
				os.Exit(1)
			}

			log.Println("Exporting to Markdown in the vault directory...")
			if err := export.ExportMarkdownFromDB(dbPath); err != nil { // Removed seedURL as it's no longer needed.
				log.Fatalf("Error exporting Markdown: %v", err)
			}
			return
		}
		if exportTo == "csv" {
			if dbPath == "" {
				fmt.Println("Database path is required for export!")
				os.Exit(1)
			}

			outputFile := dbPath + ".csv"
			log.Println("Exporting data to CSV file:", outputFile)

			export.ExportToCSV(dbPath, outputFile)
			return
		}

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
	rootCmd.Flags().StringVar(&seedURL, "seed", "", "Seed URL for crawling (required if not exporting to CSV or Markdown)")
	rootCmd.Flags().IntVar(&depth, "depth", 1, "Depth of the crawl (recursively)")
	rootCmd.Flags().StringVar(&dbPath, "db", "graph.db", "Path to the SQLite database")
	rootCmd.Flags().StringVar(&exportTo, "export", "", "Export data to a file (csv or markdown)")

	// Customize the help message
	rootCmd.SetHelpTemplate(`Usage: wombat [flags]

Flags:
{{.Flags.FlagUsages | trimTrailingWhitespaces}}

Examples:
wombat --seed "https://en.wikipedia.org/wiki/Wombat" --depth 2
wombat --export=csv --db pathto.db
wombat --export=md --db pathto.db
wombat --help
`)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
