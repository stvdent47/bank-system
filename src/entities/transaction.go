package entities

type Transaction struct {
	ID            string `db: id json: id`
	Amount        int64  `db: amount json: amount`
	FromAccountId string `db: from_id json: from`
	ToAccountId   string `db: to_id json: to`
	Type          string `db: type json: type`
	Description   string `db: description json: description`
	CreatedAt     int64  `db: created_at json: created_at`
}
