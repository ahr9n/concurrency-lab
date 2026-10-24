//go:build bughunt

package transfer

import (
	"errors"
	"sync"
)

var ErrInsufficientFunds = errors.New("insufficient funds")

type Account struct {
	ID      int
	mu      sync.Mutex
	balance int
}

func NewAccount(id, balance int) *Account {
	return &Account{ID: id, balance: balance}
}

func (a *Account) Balance() int {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.balance
}

func Transfer(from, to *Account, amount int) error {
	from.mu.Lock()
	defer from.mu.Unlock()
	to.mu.Lock()
	defer to.mu.Unlock()

	if from.balance < amount {
		return ErrInsufficientFunds
	}
	from.balance -= amount
	to.balance += amount
	return nil
}
