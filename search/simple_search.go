package search

import (
	"context"
	"fmt"
	"os"
)

func (s Search) SimpleSearch() string {
	queryUsernameString := fmt.Sprintf("SELECT username from simple where username=$1;")

	rows, err := s.db.Query(context.Background(), queryUsernameString, s.text)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to search table: %v\n", err)
		os.Exit(1)
	}

	if rows.Next() {
		var username string
		err = rows.Scan(&username)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Search results parsing failed %v", err)
			os.Exit(1)
		}
		return username
	}
	return ""
}
