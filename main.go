package main

import (
	"context"
	"ent_test/ent"
	"fmt"
	"log"
	"math"

	"entgo.io/ent/dialect/sql"
	"github.com/go-sql-driver/mysql"
)

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
	drv, err := sql.Open("mysql", mc.FormatDSN())
	client := ent.NewClient(ent.Driver(drv))
	defer client.Close()

	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}
	ctx := context.Background()
	if err := client.Schema.Create(ctx); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	_, err = client.User.Create().SetAge(math.MaxUint64).Save(ctx)
	if err != nil {
		panic(err)
	}

	rows, err := drv.QueryContext(ctx, "SELECT * FROM users")
	if err != nil {
		panic(err)
	}
	if rows.Next() {
		var value uint64
		var id int
		err := rows.Scan(&id, &value)
		if err != nil {
			panic(err)
		}
		fmt.Println(id, value)
	}
	rows.Close()

	all, err := client.User.Query().All(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println(all[0])

}
