package entities

type Card struct {
	ID         string `db:id json:id`
	Number     []byte `db:number json:number`
	Expiration []byte `db:expiration json:expiration`
	CVV        string `db:cvv json:cvv`
	UserId     string `db:user_id json:userId`
	AccountID  string `db:account_id json:accountId`
	CreatedAt  int64  `db:created_at json:createdAt`
}

type CreateCardDto struct {
	AccountID string `json:accountId`
	PGPKey    string `json:pgpKey`
}

func (this *CreateCardDto) IsValid() bool {
	if this.AccountID == "" {
		return false
	}
	if this.PGPKey == "" {
		return false
	}

	return true
}

type CreateCardResponseDto struct {
	ID         string `json:id`
	Number     string `json:number`
	Expiration string `json:expiration`
	CVV        string `json:cvv`
}

type GetCardInfoDto struct {
	CardId string `json:cardId`
	PGPKey string `json:pgpKey`
}

func (this *GetCardInfoDto) IsValid() bool {
	if this.CardId == "" {
		return false
	}
	if this.PGPKey == "" {
		return false
	}

	return true
}

type GetCardInfoResponseDto struct {
	ID         string `json:id`
	Number     string `json:number`
	Expiration string `json:expiration`
	UserId     string `json:userId`
	AccountID  string `json:accountId`
	CreatedAt  int64  `json:createdAt`
}

type PayCardDto struct {
	// todo: use card number instead:
	// CardNumber string `json:cardNumber`
	CardId     string `json:cardId`
	Expiration string `json:expiration`
	CVV        string `json:cvv`
	Amount     int64  `json:amount`
	PGPKey     string `json:pgpKey`
}

func (this *PayCardDto) IsValid() bool {
	if this.CardId == "" {
		return false
	}
	if this.Expiration == "" {
		return false
	}
	if this.CVV == "" {
		return false
	}
	if this.Amount <= 0 {
		return false
	}
	if this.PGPKey == "" {
		return false
	}

	return true
}
