package types

type WalletData struct {
	Key string
	Balance int64
	Vouchers string
}

type Voucher struct {
	Key string
	ID string
	Amount int64
	Description string
}

type Pair struct {
	Key string `json:"key"`
	Value string `json:"value"`
}