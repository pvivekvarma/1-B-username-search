package seed

import (
	"github.com/jackc/pgx/v5"
)

type SeedStrategy interface {
	Execute() error
}

type UsernamePKSeedStrategy struct {
	Db *pgx.Conn
}

func (s *UsernamePKSeedStrategy) Execute() error {
	return UsernamePKSeed(s)
}
