package model

import (
	"fmt"
	"net/http"
)

// User is a user structure.
type User struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

// Bind validates user structure.
func (u *User) Bind(r *http.Request) error {
	if u.UserName == "" {
		return fmt.Errorf("login is a required field")
	}
	if u.Password == "" {
		return fmt.Errorf("password is a required field")
	}
	return nil
}

// AuthResponse contains Auth hanlder response.
type AuthResponse struct {
	Token  string `json:"token,omitempty"`
	Errors string `json:"errors,omitempty"`
}

// Render tunes rendering of AuthResponse structure.
func (ar *AuthResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// Info is a structure, that contains information about users
// coins, inventory and transaction history.
type Info struct {
	Coins        int             `json:"coins"`
	Inventory    []InventoryItem `json:"inventory,omitempty"`
	CoinsHistory CoinsHistory    `json:"coinHistory,omitempty"`
}

// Render tunes rendering of AuthResponse structure.
func (i *Info) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// InventoryItem is an inventory item structure.
type InventoryItem struct {
	Type     string `json:"type"`
	Quantity int    `json:"quantity"`
}

// CoinsHistory contains users coin transaction history.
type CoinsHistory struct {
	Received []CoinsReceiving `json:"received,omitempty"`
	Sent     []CoinsSending   `json:"sent,omitempty"`
}

// CoinsReceiving is a coins receiving structure.
type CoinsReceiving struct {
	FromUser string `json:"fromUser"`
	Amount   int    `json:"amount"`
}

// CoinsSending is a coins sending structure.
type CoinsSending struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}

// Bind validates user structure.
func (cs *CoinsSending) Bind(r *http.Request) error {
	if cs.ToUser == "" {
		return fmt.Errorf("toUser is a required field")
	}
	if cs.Amount == 0 {
		return fmt.Errorf("amount is a required field")
	}
	return nil
}
