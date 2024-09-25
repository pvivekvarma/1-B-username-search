package seed

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func UsernamePKSeed(s *UsernamePKSeedStrategy) error {
	absPathUsernames, _ := filepath.Abs("../data/xato-net-10-million-usernames.txt")
	absPathFamilynames, _ := filepath.Abs("../data/familynames.txt")
	readUsernamesFile, err := os.Open(absPathUsernames)

	TableName := "testsimple"
	if err != nil {
		return errors.Join(err, errors.New("failed to open file"))
	}

	readFamilynamesFile, err := os.Open(absPathFamilynames)
	if err != nil {
		return errors.Join(err, errors.New("failed to open file"))
	}

	defer func(readFile *os.File) {
		err := readFile.Close()
		if err != nil {
			log.Fatalf("Exception when closing file: %v\n", err)
		}
	}(readUsernamesFile)

	defer func(readFile *os.File) {
		err := readFile.Close()
		if err != nil {
			log.Fatalf("Exception when closing file: %v\n", err)
		}
	}(readFamilynamesFile)

	dropTableString := fmt.Sprintf("DROP TABLE IF EXISTS %s;", TableName)
	createTableString := fmt.Sprintf("CREATE TABLE %s (username varchar(255) PRIMARY KEY);", TableName)

	_, err = s.Db.Exec(context.Background(), dropTableString)
	if err != nil {
		return errors.Join(err, errors.New("failed to drop table"))
	}

	_, err = s.Db.Exec(context.Background(), createTableString)
	if err != nil {
		return errors.Join(err, errors.New("failed to create table"))
		os.Exit(1)
	}

	start := time.Now()

	familyNamesFileScanner := bufio.NewScanner(readFamilynamesFile)
	familyNamesFileScanner.Split(bufio.ScanLines)
	counter := 0
	avgCounter := 0
	batchSize := 1000

	for familyNamesFileScanner.Scan() {
		familyName := familyNamesFileScanner.Text()
		fmt.Printf("FAMILYNAME: %v", familyName)
		batchStart := time.Now()
		usernames := make([]any, 0, batchSize)
		valuesPlaceholder := make([]string, 0, batchSize)

		_, err = readUsernamesFile.Seek(0, io.SeekStart)
		if err != nil {
			log.Fatal(err)
		}

		var usernamesFileScanner *bufio.Scanner = nil
		usernamesFileScanner = bufio.NewScanner(readUsernamesFile)
		usernamesFileScanner.Split(bufio.ScanLines)

		batchIndex := 0
		for usernamesFileScanner.Scan() {
			counter++
			batchIndex++
			username := fmt.Sprintf("%v%v", usernamesFileScanner.Text(), strings.ToLower(familyName))
			usernames = append(usernames, username)
			valuesPlaceholder = append(valuesPlaceholder, fmt.Sprintf("($%d)", batchIndex%batchSize+1))

			if batchIndex%batchSize == 0 {
				avgCounter++
				insertUsernameString := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s;", TableName, "username", strings.Join(valuesPlaceholder, ", "))
				_, err := s.Db.Exec(context.Background(), insertUsernameString, usernames...)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Failed to insert username %v due to %v", username, err.Error())
					if !strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
						return errors.Join(err, errors.New("failed to drop table"))
					}
				}
				fmt.Printf("Inserted %d records\n", counter)
				if avgCounter%batchSize == 0 {
					fmt.Printf("Took average of %v to insert %d records into database\n", time.Now().Sub(batchStart)/1000, batchSize)
					batchStart = time.Now()
				}
				usernames = nil
				valuesPlaceholder = nil
			}
		}
	}
	fmt.Printf("Took total of %v to insert %d records into database\n", time.Now().Sub(start), counter)

	return nil
}
