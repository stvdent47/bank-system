package entities

type Loan struct {
	ID           string  `db:id json:id`
	UserId       string  `db:user_id json:userId`
	AccountId    string  `db:account_id json:accountId`
	Amount       int64   `db:amount json:amount`
	InterestRate float64 `db:interest_rate json:interestRate`
	Term         int     `db:term json:term`
	StartDate    int64   `db:start_date json:startDate`
	Debt         int64   `db:debt json:debt`
}

type LoanApplyDto struct {
	UserId    string `json:userId`
	AccountId string `json:accountId`
	Amount    int64  `json:amount`
	Term      int    `json:term`
}

func (this *LoanApplyDto) IsValid() bool {
	if this.UserId == "" {
		return false
	}
	if this.AccountId == "" {
		return false
	}
	if this.Amount <= 0 {
		return false
	}
	if this.Term <= 0 {
		return false
	}

	return true
}

type LoanResponseDto struct {
	ID           string    `json:id`
	UserId       string    `json:userId`
	AccountId    string    `json:accountId`
	Amount       int64     `json:amount`
	InterestRate float64   `json:interestRate`
	Term         int       `json:term`
	StartDate    int64     `json:startDate`
	Debt         int64     `json:debt`
	Payments     []Payment `json:payments`
}
