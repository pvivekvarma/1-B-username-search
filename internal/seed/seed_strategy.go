package seed

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type SeedStrategy interface {
	Execute() error
}

type UsernamePKSeedStrategy struct {
	Db *pgxpool.Pool
}

func (s *UsernamePKSeedStrategy) Execute() error {
	return UsernamePKSeed(s)
}


type UsernameSeedStrategy struct {
	Db *pgxpool.Pool
}

func (s *UsernameSeedStrategy) Execute() error {
	return UsernameSeed(s)
}
