package entities

type Payment struct {
	ID            string `db:id json:id`
	LoanId        string `db:loan_id json:loanId`
	Amount        int64  `db:amount json:amount`
	Date          int64  `db:date json:date`
	DueDate       int64  `db:due_date json:dueDate`
	PrincipalPart int64  `db:principal_part json:principalPart`
	InterestPart  int64  `db:interest_part json:interestPart`
	Status        string `db:status json:status`
	IsPaid        bool   `db:is_paid json:isPaid`
}
