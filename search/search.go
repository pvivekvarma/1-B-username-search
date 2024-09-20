package search

import (
	"github.com/jackc/pgx/v5"
	"log"
	"reflect"
	"time"
)

type Search struct {
	db   *pgx.Conn
	text string
}

func Execute(db *pgx.Conn, searchMethodName, searchString string) (string, time.Duration) {
	s := Search{db, searchString}
	begin := time.Now()
	return search(s, searchMethodName), time.Now().Sub(begin)
}

func search(s Search, searchMethodName string) string {
	searchMethod := reflect.ValueOf(s).MethodByName(searchMethodName)

	if !searchMethod.IsValid() {
		log.Fatal("Invalid seed method given")
	}

	log.Printf("Searching username using method %v ...\n", searchMethodName)
	val := searchMethod.Call(nil)
	return val[0].String()
}
