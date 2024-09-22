package seeds

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func (s Seed) PgSimpleUsernameSeed() {
	absPathUsernames, _ := filepath.Abs("./data/xato-net-10-million-usernames.txt")
	absPathFamilynames, _ := filepath.Abs("./data/familynames.txt")
	readUsernamesFile, err := os.Open(absPathUsernames)
	if err != nil {
		fmt.Println(err)
	}
	readFamilynamesFile, err := os.Open(absPathFamilynames)
	if err != nil {
		fmt.Println(err)
	}

	defer func(readFile *os.File) {
		err := readFile.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Exception when closing file: %v\n", err)
			os.Exit(1)
		}
	}(readUsernamesFile)

	defer func(readFile *os.File) {
		err := readFile.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Exception when closing file: %v\n", err)
			os.Exit(1)
		}
	}(readFamilynamesFile)

	dropTableString := fmt.Sprintf("DROP TABLE IF EXISTS simple;")
	createTableString := fmt.Sprintf("CREATE TABLE simple (username varchar(255) PRIMARY KEY);")

	_, err = s.db.Exec(context.Background(), dropTableString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to drop table: %v\n", err)
		os.Exit(1)
	}

	_, err = s.db.Exec(context.Background(), createTableString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to drop table: %v\n", err)
		os.Exit(1)
	}

	start := time.Now()

	familynamesFileScanner := bufio.NewScanner(readFamilynamesFile)
	familynamesFileScanner.Split(bufio.ScanLines)
	counter := 0
	avgCounter := 0
	batchSize := 1000

	for familynamesFileScanner.Scan() {
		familyname := familynamesFileScanner.Text()
		fmt.Printf("FAMILYNAME: %v", familyname)
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
			username := fmt.Sprintf("%v%v", usernamesFileScanner.Text(), strings.ToLower(familyname))
			usernames = append(usernames, username)
			valuesPlaceholder = append(valuesPlaceholder, fmt.Sprintf("($%d)", batchIndex%batchSize+1))

			if batchIndex%batchSize == 0 {
				avgCounter++
				insertUsernameString := fmt.Sprintf("INSERT INTO simple (%s) VALUES %s;", "username", strings.Join(valuesPlaceholder, ", "))
				_, err := s.db.Exec(context.Background(), insertUsernameString, usernames...)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Failed to insert username %v due to %v", username, err.Error())
					if !strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
						os.Exit(1)
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
}
