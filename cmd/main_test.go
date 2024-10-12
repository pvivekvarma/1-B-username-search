package main

import (
	"com/pvivekvarma/1-B-username-search/internal/search"
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var (
	randUsernames []string

	conn *pgxpool.Pool
)

func load(b *testing.B) {
	err := godotenv.Load()

	connString := fmt.Sprintf("postgresql://%v:%v@%v:%v/%v", os.Getenv("PG_USERNAME"), os.Getenv("PG_PASSWORD"), os.Getenv("PG_HOST"), os.Getenv("PG_PORT"), os.Getenv("PG_NAME"))
	conn, err = pgxpool.New(context.Background(), connString)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	randUsernames = make([]string, 0)
	queryRandomUsernames := fmt.Sprintf("SELECT username from testsimple ORDER BY random() LIMIT %d", 100)
	rows, err := conn.Query(context.Background(), queryRandomUsernames)
	if err != nil {
		b.Fatalf("Failed querying random uesrnames: %v", err.Error())
	}

	for rows.Next() {
		var username string
		err := rows.Scan(&username)
		if err != nil {
			b.Fatal("Failed parsing results while loading")
		}
		randUsernames = append(randUsernames, username)
	}
	rows.Close()

	fmt.Printf("%v", randUsernames)
}

func BenchmarkUsernamePKSearch(b *testing.B) {
	load(b)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		usernamePkSearch := &search.UsernamePKSearchStrategy{
			Db:         conn,
			SearchText: randUsernames[i%len(randUsernames)],
		}

		err := usernamePkSearch.Execute()
		if err != nil {
			b.Fatalf("Benchmark failed %v", err)
		}
	}
}
