package spirit

import (
	"fmt"
	"log"

	"github.com/tenmozes/sanchez/pkg/pump"
)

type Service struct {
	pumps     map[string]pump.Worker
	cocktails map[string]Cocktail
}

func NewService(pumps []pump.Worker, cocktails []Cocktail) (*Service, error) {
	s := &Service{
		pumps:     make(map[string]pump.Worker, len(pumps)),
		cocktails: make(map[string]Cocktail, len(cocktails)),
	}
	for i := range pumps {
		s.pumps[pumps[i].Name()] = pumps[i]
	}
	for i, c := range cocktails {
		if len(cocktails[i].Ingredient) == 0 {
			return nil, fmt.Errorf("no ingredient in coctails %s", c.Name)
		}
		s.cocktails[c.Name] = c

	}
	return s, nil
}

func (s *Service) Make(name string) error {
	cocktail, ok := s.cocktails[name]
	if !ok {
		return CocktailError{name: name}
	}
	for _, ingrdnt := range cocktail.Ingredient {
		worker, ok := s.pumps[ingrdnt.Type]
		if !ok {
			return WorkerError{name: ingrdnt.Type}
		}
		if err := worker.Work(ingrdnt.Duration); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) Close() error {
	for i, p := range s.pumps {
		if err := p.Close(); err != nil {
			log.Printf("warning: can't close worker %s, %v", i, err)
		}
	}
	return nil
}
