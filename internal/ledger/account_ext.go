package ledger

// balances returns a deep copy of the account's current balance map.
func (a *Account) balances() map[string]float64 {
	a.mu.RLock()
	defer a.mu.RUnlock()
	copy := make(map[string]float64, len(a.balance))
	for cur, val := range a.balance {
		copy[cur] = val
	}
	return copy
}

// CreditTotal returns the sum of all credit entries posted to the account
// across all currencies converted to a single float for reporting purposes.
// Note: this is a raw sum and does not perform FX conversion.
func (a *Account) CreditTotal() float64 {
	a.mu.RLock()
	defer a.mu.RUnlock()
	var total float64
	for _, v := range a.credits {
		total += v
	}
	return total
}

// DebitTotal returns the sum of all debit entries posted to the account
// across all currencies.
func (a *Account) DebitTotal() float64 {
	a.mu.RLock()
	defer a.mu.RUnlock()
	var total float64
	for _, v := range a.debits {
		total += v
	}
	return total
}
