package seeds

import (
	"github.com/jackc/pgx/v5"
	"log"
	"reflect"
)

type Seed struct {
	db *pgx.Conn
}

func Execute(db *pgx.Conn, seedMethodName string) {
	s := Seed{db}
	log.Printf("Seeding method %v", seedMethodName)
	seed(s, seedMethodName)
}

func seed(s Seed, seedMethodName string) {
	seedMethod := reflect.ValueOf(s).MethodByName(seedMethodName)

	if !seedMethod.IsValid() {
		log.Fatal("Invalid seed method given")
	}

	log.Printf("Seeding database using method %v ...\n", seedMethodName)
	seedMethod.Call(nil)
	log.Printf("Completed seeding database using method %v\n", seedMethodName)
}
