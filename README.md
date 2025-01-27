# Wombat 🐨

> Wombat is a Wikipedia crawler built in Go that helps visualize relationships between links in Wikipedia articles.

```sh
go run github.com/polarhive/wombat@latest --seed "https://en.wikipedia.org/wiki/Wombat"
```

## Usage

> You can specify additional flags to control the depth and database path.

### Flags

- **--seed**: The seed URL to start crawling from (required).
- **--depth**: The depth of the crawl. Increasing depth crawls more links.
- **--db**: Path to the SQLite database (default is `graph.db`).
- **--export**: Export the crawled data to a file in CSV or Markdown format (Obsidian).

### Example

```sh
wombat --seed "https://en.wikipedia.org/wiki/Wombat" --depth 2
wombat --seed "https://en.wikipedia.org/wiki/Wombat" --db /path/to/custom.db
wombat --export=csv --db /path/to/db
wombat --export=md --db /path/to/db
```
---

For more information, visit: [polarhive.net/wombat](https://github.com/polarhive/wombat)
