# ledger-sync

A lightweight double-entry ledger reconciliation tool for microservices.

---

## Installation

```bash
go install github.com/your-org/ledger-sync@latest
```

Or clone and build manually:

```bash
git clone https://github.com/your-org/ledger-sync.git
cd ledger-sync && go build ./...
```

---

## Usage

Define your ledger entries and run a reconciliation check:

```go
package main

import "github.com/your-org/ledger-sync/ledger"

func main() {
    l := ledger.New()

    l.Post(ledger.Entry{
        Debit:  "accounts_receivable",
        Credit: "revenue",
        Amount: 1000_00, // in cents
        Ref:    "inv-001",
    })

    if err := l.Reconcile(); err != nil {
        log.Fatalf("reconciliation failed: %v", err)
    }

    log.Println("Ledger balanced ✓")
}
```

Run from the CLI:

```bash
ledger-sync reconcile --config ledger.yaml
```

---

## Features

- Double-entry validation across distributed services
- Pluggable storage backends (Postgres, Redis, in-memory)
- CLI and library interfaces
- Lightweight with zero external dependencies in core package

---

## Contributing

Pull requests are welcome. Please open an issue first to discuss any significant changes.

---

## License

[MIT](LICENSE) © your-org