package repository

const createSubscription = `INSERT INTO sh_eff.subscriptions
	(id, name, price, user_id, start_date, end_date, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
	RETURNING id, name, price, user_id, start_date, end_date, created_at, updated_at`
