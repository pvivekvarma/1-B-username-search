package search

import "com/pvivekvarma/1-B-username-search/internal/command"

type SearchCommand struct {
	Strategy SearchStrategy
	next     command.Command
}

func (s *SearchCommand) Execute() error {
	err := s.Strategy.Execute()
	if s.next != nil && err != nil {
		err = s.next.Execute()
	}
	return err
}

func (s *SearchCommand) SetNext(c command.Command) {
	s.next = c
}
