package ledger

// Currencies returns the list of currencies for which the account has
// recorded at least one entry. The order is not guaranteed.
func (a *Account) Currencies() []string {
	seen := make(map[string]struct{})
	for _, e := range a.entries {
		seen[e.Currency] = struct{}{}
	}
	out := make([]string, 0, len(seen))
	for c := range seen {
		out = append(out, c)
	}
	return out
}

// Balance returns the net balance for the given currency.
// For asset/expense accounts a debit increases the balance;
// for liability/equity/revenue accounts a credit increases it.
func (a *Account) Balance(currency string) float64 {
	var balance float64
	for _, e := range a.entries {
		if e.Currency != currency {
			continue
		}
		switch a.Type {
		case AccountTypeAsset, AccountTypeExpense:
			if e.Type == EntryTypeDebit {
				balance += e.Amount
			} else {
				balance -= e.Amount
			}
		default:
			if e.Type == EntryTypeCredit {
				balance += e.Amount
			} else {
				balance -= e.Amount
			}
		}
	}
	return balance
}
