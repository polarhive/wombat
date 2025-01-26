# Wombat \\/(^.^)\/\

> Wombat is a Wikipedia crawler built in Go that helps visualize relationships between links in Wikipedia articles.

```sh
go run github.com/polarhive/wombat --seed "https://en.wikipedia.org/wiki/Wombat"
```

## Usage

You can specify additional flags to control the depth and database path.

- **--seed**: The seed URL to start crawling from (required).
- **--depth**: The depth of the crawl (default is 1).
- **--db**: Path to the SQLite database (default is `graph.db`).


## Benchmarks

You can run a benchmark to test the performance of the crawler with the following command:

```
go test -bench .
```

**Sample Benchmark Output:**
```text
todo
```


---

For more information, visit: [polarhive.net/wombat](https://github.com/polarhive/wombat)
