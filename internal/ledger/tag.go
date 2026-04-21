package ledger

import (
	"errors"
	"sort"
	"strings"
)

// Tag represents a key-value metadata label attached to a transaction.
type Tag struct {
	Key   string
	Value string
}

// TagIndex maintains a mapping from tag keys/values to transaction IDs.
type TagIndex struct {
	index map[string][]string // "key:value" -> []transactionID
}

// NewTagIndex creates an empty TagIndex.
func NewTagIndex() *TagIndex {
	return &TagIndex{
		index: make(map[string][]string),
	}
}

// Add indexes a transaction ID under the given tag.
func (ti *TagIndex) Add(txID string, tag Tag) error {
	if strings.TrimSpace(tag.Key) == "" {
		return errors.New("tag key must not be empty")
	}
	if strings.TrimSpace(txID) == "" {
		return errors.New("transaction ID must not be empty")
	}
	composite := tag.Key + ":" + tag.Value
	ti.index[composite] = append(ti.index[composite], txID)
	return nil
}

// Lookup returns all transaction IDs associated with the given tag.
func (ti *TagIndex) Lookup(tag Tag) []string {
	composite := tag.Key + ":" + tag.Value
	ids, ok := ti.index[composite]
	if !ok {
		return []string{}
	}
	// Return a sorted copy to ensure deterministic output.
	copy_ := make([]string, len(ids))
	copy(copy_, ids)
	sort.Strings(copy_)
	return copy_
}

// Keys returns all unique composite "key:value" strings stored in the index.
func (ti *TagIndex) Keys() []string {
	keys := make([]string, 0, len(ti.index))
	for k := range ti.index {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// Remove deletes all entries for the given tag from the index.
func (ti *TagIndex) Remove(tag Tag) {
	composite := tag.Key + ":" + tag.Value
	delete(ti.index, composite)
}
