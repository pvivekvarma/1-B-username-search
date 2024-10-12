package search

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SearchStrategy interface {
	Execute() error
}

type UsernamePKSearchStrategy struct {
	Db         *pgxpool.Pool
	SearchText string
}

func (s *UsernamePKSearchStrategy) Execute() error {
	fmt.Print("Searching using UsernamePKSearchStrategy\n")
	TableName := "testsimple"

	queryUsernameString := fmt.Sprintf("SELECT username from %s where username=$1;", TableName)

	rows, err := s.Db.Query(context.Background(), queryUsernameString, s.SearchText)
	if err != nil {
		return errors.Join(err, errors.New("failed to search table"))
	}

	if rows.Next() {
		var username string
		err = rows.Scan(&username)
		if err != nil {
			return errors.Join(err, errors.New("search results parsing failed"))
		}
		if username == s.SearchText {
			fmt.Printf("Username %s found!\n", s.SearchText)
		} else {
			fmt.Printf("Username %s not found!\n", s.SearchText)
		}
	}

	rows.Close()

	return nil
}
