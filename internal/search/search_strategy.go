package search

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log"
	"os"
)

type SearchStrategy interface {
	Execute() error
}

type SimpleUsernamePKSearchStrategy struct {
	Db         *pgx.Conn
	SearchText string
}

func (s *SimpleUsernamePKSearchStrategy) Execute() error {
	TableName := "testsimple"

	queryUsernameString := fmt.Sprintf("SELECT username from %s where username=$1;", TableName)

	rows, err := s.Db.Query(context.Background(), queryUsernameString, s.SearchText)
	if err != nil {
		return err
		fmt.Fprintf(os.Stderr, "Failed to search table: %v\n", err)
		os.Exit(1)
	}

	if rows.Next() {
		var username string
		err = rows.Scan(&username)
		if err != nil {
			return err
			log.Fatalf("Search results parsing failed %v", err)
		}
		return err
	}
	return err
}
