package ledger

import "fmt"

// RuleViolation describes a single rule breach.
type RuleViolation struct {
	RuleID    string
	AccountID string
	Message   string
}

// CheckRules evaluates all rules in ri against book b.
// It returns a slice of violations (empty means compliant).
func CheckRules(b *Book, ri *RuleIndex, tags *TagIndex) ([]RuleViolation, error) {
	if b == nil {
		return nil, fmt.Errorf("check_rules: book must not be nil")
	}
	if ri == nil {
		return nil, fmt.Errorf("check_rules: rule index must not be nil")
	}

	var violations []RuleViolation

	for _, rule := range ri.All() {
		acc, err := b.GetAccount(rule.AccountID)
		if err != nil {
			// Account not present — skip rule
			continue
		}

		switch rule.Type {
		case RuleTypeMinBalance:
			bal := acc.Balance(rule.Currency)
			if bal < rule.Threshold {
				violations = append(violations, RuleViolation{
					RuleID:    rule.ID,
					AccountID: rule.AccountID,
					Message: fmt.Sprintf("balance %d %s below minimum %d",
						bal, rule.Currency, rule.Threshold),
				})
			}

		case RuleTypeMaxBalance:
			bal := acc.Balance(rule.Currency)
			if bal > rule.Threshold {
				violations = append(violations, RuleViolation{
					RuleID:    rule.ID,
					AccountID: rule.AccountID,
					Message: fmt.Sprintf("balance %d %s exceeds maximum %d",
						bal, rule.Currency, rule.Threshold),
				})
			}

		case RuleTypeRequireTag:
			if tags == nil {
				break
			}
			txIDs := tags.Lookup(rule.TagKey)
			if len(txIDs) == 0 {
				violations = append(violations, RuleViolation{
					RuleID:    rule.ID,
					AccountID: rule.AccountID,
					Message: fmt.Sprintf("no transactions tagged with key %q", rule.TagKey),
				})
			}
		}
	}

	return violations, nil
}

// HasRuleViolations returns true when CheckRules finds at least one violation.
func HasRuleViolations(b *Book, ri *RuleIndex, tags *TagIndex) bool {
	v, _ := CheckRules(b, ri, tags)
	return len(v) > 0
}
