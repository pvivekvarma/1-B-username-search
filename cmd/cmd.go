package main

import (
	"com/pvivekvarma/1-B-username-search/internal/command"
	"com/pvivekvarma/1-B-username-search/internal/search"
	"com/pvivekvarma/1-B-username-search/internal/seed"
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var (
	isSeed         bool
	searchStrategy string
	searchText     string
)

var SearchTypes = [...]string{"pg_username_pk", "pg_username"}

func main() {
	fmt.Println("Search 1 billion usernames (100 million in this case)")
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}

	handleArgs()
}

func handleArgs() {
	flag.BoolVar(&isSeed, "seed", false, "Seed database?")
	flag.StringVar(&searchStrategy, "strategy", "pg_username_pk", "Type of search")
	flag.StringVar(&searchText, "search", "", "The username to search")

	flag.Parse()

	var isValidSearch = false
	for _, s := range SearchTypes {
		if s == searchStrategy {
			isValidSearch = true
		}
	}

	if !isValidSearch {
		log.Fatal("The given search type is invalid.")
	}

	switch searchStrategy {
	case "pg_username_pk":
		connString := fmt.Sprintf("postgresql://%v:%v@%v:%v/%v", os.Getenv("PG_USERNAME"), os.Getenv("PG_PASSWORD"), os.Getenv("PG_HOST"), os.Getenv("PG_PORT"), os.Getenv("PG_NAME"))
		conn, err := pgxpool.New(context.Background(), connString)
		if err != nil {
			log.Fatalf("Unable to connect to database: %v\n", err)
		}

		var c command.Command
		if isValidSearch {
			search := &search.SearchCommand{
				Strategy: &search.UsernamePKSearchStrategy{
					Db:         conn,
					SearchText: searchText,
				},
			}
			c = search
		}

		if isSeed {
			seed := &seed.SeedCommand{
				Strategy: &seed.UsernamePKSeedStrategy{
					Db: conn,
				},
				Seed: isSeed,
			}
			if c != nil {
				seed.SetNext(c)
				c = seed
			}
		}

		defer conn.Close()
		if c != nil {
			err = c.Execute()
			if err != nil {
				log.Fatalf("Program failed with a non-zero exit code: %v\n", err)
			}
		}

	case "pg_username":
		connString := fmt.Sprintf("postgresql://%v:%v@%v:%v/%v", os.Getenv("PG_USERNAME"), os.Getenv("PG_PASSWORD"), os.Getenv("PG_HOST"), os.Getenv("PG_PORT"), os.Getenv("PG_NAME"))
		conn, err := pgxpool.New(context.Background(), connString)
		if err != nil {
			log.Fatalf("Unable to connect to database: %v\n", err)
		}

		var c command.Command
		if isValidSearch {
			search := &search.SearchCommand{
				Strategy: &search.UsernameSearchStrategy{
					Db:         conn,
					SearchText: searchText,
				},
			}
			c = search
		}

		if isSeed {
			seed := &seed.SeedCommand{
				Strategy: &seed.UsernameSeedStrategy{
					Db: conn,
				},
				Seed: isSeed,
			}
			if c != nil {
				seed.SetNext(c)
				c = seed
			}
		}

		defer conn.Close()
		if c != nil {
			err = c.Execute()
			if err != nil {
				log.Fatalf("Program failed with a non-zero exit code: %v\n", err)
			}
		}
	}
}
