package seed

import "com/pvivekvarma/1-B-username-search/internal/command"

type SeedCommand struct {
	Strategy SeedStrategy
	Seed     bool
	next     command.Command
}

func (s *SeedCommand) Execute() error {
	var err error
	if s.Seed == true {
		err = s.Strategy.Execute()
	}
	if s.next != nil && err != nil {
		err = s.next.Execute()
	}
	return err
}

func (s *SeedCommand) SetNext(c command.Command) {
	s.next = c
}
