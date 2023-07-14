package main

import (
	"context"
	"ent_test/ent"
	"log"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/schema"
	"github.com/go-sql-driver/mysql"
)

var client *ent.Client

func main() {
	mc := mysql.Config{
		User:                 "root",
		Passwd:               "",
		Net:                  "tcp",
		Addr:                 "localhost" + ":" + "3306",
		DBName:               "ent",
		AllowNativePasswords: true,
		ParseTime:            true,
	}
	drv, _ := sql.Open("mysql", mc.FormatDSN())
	client = ent.NewClient(ent.Driver(drv), ent.Debug())
	ctx := context.Background()
	err := client.Schema.Create(ctx, schema.WithDropIndex(true), schema.WithDropColumn(true), schema.WithForeignKeys(true))
	if err != nil {
		log.Fatal(err)
		return
	}
	client.User.Delete().ExecX(ctx)
	client.Group.Delete().ExecX(ctx)
	first()

	second()

}

func first() {
	ctx := context.Background()

	// Begin transaction
	tx, err := client.Tx(ctx)
	if err != nil {
		log.Fatalf("failed creating transaction: %v", err)
	}

	users := []*ent.User{
		{Name: "Alice", Age: 20},
		{Name: "Bob", Age: 25},
	}

	for _, u := range users {
		_, err = tx.User.Create().SetName(u.Name).SetAge(u.Age).Save(ctx)
		if err != nil {
			log.Fatalf("failed creating user: %v", err)
		}
	}

	if err := tx.Commit(); err != nil {
		log.Fatalf("failed committing transaction: %v", err)
	}
}

func second() {
	ctx := context.Background()

	tx, err := client.Tx(ctx)
	if err != nil {
		log.Fatalf("failed creating transaction: %v", err)
	}

	newUsers := []*ent.User{
		{Name: "Charlie", Age: 30},
		{Name: "Bob", Age: 27},
	}

	groups := []*ent.Group{
		{Name: "Group1", Value: "value1"},
		{Name: "Group2", Value: "value2"},
	}

	entGroups := make([]*ent.Group, len(groups))
	for i, g := range groups {
		entGroups[i], err = tx.Group.Create().SetName(g.Name).SetValue(g.Value).Save(ctx)
		if err != nil {
			log.Fatalf("failed creating group: %v", err)
		}
	}

	bulk := make([]*ent.UserCreate, len(newUsers))
	for i, u := range newUsers {
		uc := client.User.Create().
			SetAge(u.Age).
			SetName(u.Name)

		for _, g := range entGroups {
			uc = uc.AddGroupIDs(g.ID)
		}

		bulk[i] = uc
	}

	err = tx.User.
		CreateBulk(bulk...).
		OnConflict().
		UpdateNewValues().
		Exec(ctx)
	if err != nil {
		log.Fatalf("failed creating users: %v", err)
	}

	if err := tx.Commit(); err != nil {
		log.Fatalf("failed committing transaction: %v", err)
	}
}
