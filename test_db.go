package main

import (
	"context"
	"fmt"
	"log"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	db, err := pgxpool.New(context.Background(), "postgres://postgres:postgres@localhost:5432/finance_app?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var id int
	err = db.QueryRow(context.Background(), "select a.id, coalesce(sum(t.amount), 0) from accounts a left join transactions t on a.id = t.account_id where a.id = 1").Scan(&id, new(float64))
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Success:", id)
	}
}
