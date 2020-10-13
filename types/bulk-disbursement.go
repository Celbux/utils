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

type Shard struct {
	Amount int64
	Filled int64
	ErrorChildren []string
	Children []string
	ID string
}

type DisburseRow struct {
	Amount string `json:"amount"`
	To string `json:"to"`
}

type Response struct {
	Data string `json:"data"`
	Error string `json:"error"`
}


type Recon struct {
	TotalDisburse int64
	TotalDisbursees int
}