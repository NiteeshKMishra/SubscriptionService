package database

import (
	"context"
	"fmt"
	"log"
	"time"
)

type Plan struct {
	ID                  string
	PlanName            string
	PlanAmount          int
	PlanAmountFormatted string
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

func (p *Plan) GetAll() ([]*Plan, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, plan_name, plan_amount, created_at, updated_at
	from plans order by id`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plans []*Plan

	for rows.Next() {
		var plan Plan
		err := rows.Scan(
			&plan.ID,
			&plan.PlanName,
			&plan.PlanAmount,
			&plan.CreatedAt,
			&plan.UpdatedAt,
		)

		plan.PlanAmountFormatted = plan.AmountForDisplay()
		if err != nil {
			log.Println("Error scanning", err)
			return nil, err
		}

		plans = append(plans, &plan)
	}

	return plans, nil
}

func (p *Plan) GetOne(id string) (*Plan, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, plan_name, plan_amount, created_at, updated_at from plans where id = $1`

	var plan Plan
	row := db.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&plan.ID,
		&plan.PlanName,
		&plan.PlanAmount,
		&plan.CreatedAt,
		&plan.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}
	plan.PlanAmountFormatted = plan.AmountForDisplay()

	return &plan, nil
}

// SubscribeUserToPlan subscribes a user to one plan by insert
// values into user_plans table inside transaction
func (p *Plan) SubscribeUserToPlan(userID string, planID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	// delete existing plan, if any
	stmt := `delete from user_plans where user_id = $1`
	_, err = tx.ExecContext(ctx, stmt, userID)
	if err != nil {
		return err
	}

	// subscribe to new plan
	stmt = `insert into user_plans (user_id, plan_id, created_at, updated_at)
			values ($1, $2, $3, $4)`

	_, err = tx.ExecContext(ctx, stmt, userID, planID, time.Now(), time.Now())
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (p *Plan) AmountForDisplay() string {
	amount := float64(p.PlanAmount) / 100.0
	return fmt.Sprintf("$%.2f", amount)
}
