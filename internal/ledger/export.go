package ledger

import (
	"encoding/csv"
	"fmt"
	"io"
	"time"
)

// ExportFormat defines the supported export formats.
type ExportFormat string

const (
	FormatCSV ExportFormat = "csv"
)

// ExportOptions configures the export behaviour.
type ExportOptions struct {
	Format    ExportFormat
	AccountID string // if non-empty, filter to a single account
}

// ExportTransactions writes all transactions in the book to w in the
// requested format. Only FormatCSV is currently supported.
func ExportTransactions(b *Book, w io.Writer, opts ExportOptions) error {
	if b == nil {
		return fmt.Errorf("export: book must not be nil")
	}
	switch opts.Format {
	case FormatCSV, "":
		return exportCSV(b, w, opts)
	default:
		return fmt.Errorf("export: unsupported format %q", opts.Format)
	}
}

func exportCSV(b *Book, w io.Writer, opts ExportOptions) error {
	cw := csv.NewWriter(w)

	header := []string{"transaction_id", "posted_at", "entry_type", "account_id", "amount", "currency", "description"}
	if err := cw.Write(header); err != nil {
		return fmt.Errorf("export csv: write header: %w", err)
	}

	b.mu.RLock()
	txns := make([]*Transaction, len(b.transactions))
	copy(txns, b.transactions)
	b.mu.RUnlock()

	for _, tx := range txns {
		for _, e := range tx.Entries {
			if opts.AccountID != "" && e.AccountID != opts.AccountID {
				continue
			}
			row := []string{
				tx.ID,
				tx.PostedAt.Format(time.RFC3339),
				string(e.Type),
				e.AccountID,
				formatAmount(e.Amount),
				e.Currency,
				tx.Description,
			}
			if err := cw.Write(row); err != nil {
				return fmt.Errorf("export csv: write row: %w", err)
			}
		}
	}

	cw.Flush()
	return cw.Error()
}

// formatAmount renders an int64 minor-unit amount as a decimal string
// with two decimal places (e.g. 1050 → "10.50").
func formatAmount(minor int64) string {
	sign := ""
	if minor < 0 {
		sign = "-"
		minor = -minor
	}
	return fmt.Sprintf("%s%d.%02d", sign, minor/100, minor%100)
}
