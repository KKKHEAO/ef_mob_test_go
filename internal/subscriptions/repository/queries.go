package repository

const (
	createSubscription = `INSERT INTO sh_eff.subscriptions
	(id, name, price, user_id, start_date, end_date, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
	RETURNING id, name, price, user_id, start_date, end_date, created_at, updated_at`

	getSubscriptionByID = `SELECT id, name, price, user_id, start_date, end_date, created_at, updated_at FROM sh_eff.subscriptions WHERE id = $1`

	deleteSubscriptionByID = `DELETE FROM sh_eff.subscriptions WHERE id = $1`

	updateSubscriptionByID = `UPDATE sh_eff.subscriptions SET name = $2, price = $3, end_date = $4, updated_at = CURRENT_TIMESTAMP
	WHERE id = $1 RETURNING id, name, price, user_id, start_date, end_date, created_at, updated_at`

	listSubscriptions = `SELECT id, name, price, user_id, start_date, end_date, created_at, updated_at
	FROM sh_eff.subscriptions
	ORDER BY created_at DESC
	LIMIT $1 OFFSET $2`

	countSubscriptions = `SELECT COUNT(*) FROM sh_eff.subscriptions`
)
