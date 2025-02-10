package model

// User is a user structure.
type User struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

// Info is a structure, that contains information about users
// coins, inventory and transaction history.
type Info struct {
	Coins        int             `json:"coins"`
	Inventory    []InventoryItem `json:"inventory"`
	CoinsHistory CoinsHistory    `json:"coinHistory"`
}

// InventoryItem is an inventory item structure.
type InventoryItem struct {
	Type     string `json:"type"`
	Quantity string `json:"quantity"`
}

// CoinsHistory contains users coin transaction history.
type CoinsHistory struct {
	Received []CoinsReceiving `json:"received"`
	Sent     []CoinsSending   `json:"sent"`
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
