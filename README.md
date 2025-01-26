# Wombat \\/(^.^)\/\

> Wombat is a Wikipedia crawler built in Go that helps visualize relationships between links in Wikipedia articles.

---

## Setup

```
go install github.com/polarhive/wombat
```


## Usage


To start crawling a Wikipedia article, use the following command:

```
wombat --seed "https://en.wikipedia.org/wiki/Wombat"
```


## Flags

- **--seed**: seed url
- **--depth**: yes
- **--output**: json? .db?

---

## Benchmarks

Here's an example benchmark for crawling a single article:

```bash
go test -bench .
```

**Sample Benchmark Output:**
```text
todo
```
