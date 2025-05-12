package entities

type Account struct {
	ID        string `db:id json:id`
	Balance   int64  `db:balance json:balance`
	UserID    string `db:user_id json:userId`
	CreatedAt int64  `db:created_at json:createdAt`
}

type UpdateAccountBalanceDto struct {
	Amount int64  `json:amount`
	Type   string `json:type`
}

func (this *UpdateAccountBalanceDto) IsValid() bool {
	if this.Amount == 0 {
		return false
	}
	if this.Type != "deposit" && this.Type != "withdrawal" && this.Type != "transfer" && this.Type != "payment" && this.Type != "loan" {
		return false
	}

	return true
}

type TransferDto struct {
	FromAccountId string `json:fromAccountId`
	ToAccountId   string `json:toAccountId`
	Amount        int64  `json:amount`
	Description   string `json:description`
}

func (this *TransferDto) IsValid() bool {
	if this.FromAccountId == "" {
		return false
	}
	if this.ToAccountId == "" {
		return false
	}
	if this.Amount <= 0 {
		return false
	}

	return true
}
