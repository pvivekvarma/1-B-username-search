package main

import (
	"com/pvivekvarma/1-B-username-search/search"
	"com/pvivekvarma/1-B-username-search/seeds"
	"context"
	"flag"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"os"
)

var (
	seed       bool
	searchType string
	searchText string
)

var SearchTypes = [...]string{"pg_simple"}

func main() {
	fmt.Println("Search 1 billion usernames (10 million in this case)")
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}

	handleArgs()
}

func handleArgs() {
	flag.BoolVar(&seed, "seed", false, "Seed database?")
	flag.StringVar(&searchType, "search_type", "pg_simple", "Type of search")
	flag.StringVar(&searchText, "search", "", "The username to search")

	flag.Parse()

	var isValidSearch = false
	for _, s := range SearchTypes {
		if s == searchType {
			isValidSearch = true
		}
	}

	if !isValidSearch {
		fmt.Println("The given search type is invalid.")
		os.Exit(1)
	}

	switch searchType {
	case "pg_simple":
		connString := fmt.Sprintf("postgresql://%v:%v@%v:%v/%v", os.Getenv("PG_USERNAME"), os.Getenv("PG_PASSWORD"), os.Getenv("PG_HOST"), os.Getenv("PG_PORT"), os.Getenv("PG_NAME"))
		conn, err := pgx.Connect(context.Background(), connString)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
			os.Exit(1)
		}
		defer conn.Close(context.Background())

		if seed {
			seeds.Execute(conn, "PgSimpleUsernameSeed")
		}
		if searchText != "" {
			searchResult, duration := search.Execute(conn, "SimpleSearch", searchText)
			fmt.Printf("Search took %v\n", duration)
			if searchResult != "" && searchResult == searchText {
				fmt.Printf("Username %v exists!\n", searchResult)
			} else {
				fmt.Printf("Username %v does not exist!\n", searchText)
			}
		}

	}
}
