package spirit

import (
	"fmt"
	"time"
)

type Cocktail struct {
	Name       string       `yaml:"name"`
	Ingredient []Ingredient `yaml:"ingredients"`
}

type Ingredient struct {
	Type     string        `yaml:"type"`
	Duration time.Duration `yaml:"duration"`
	Order    int           `yaml:"order"`
}

type CocktailError struct {
	name string
}

func (ce CocktailError) Error() string {
	return fmt.Sprintf("%s is not exist", ce.name)
}

type WorkerError struct {
	name string
}

func (we WorkerError) Error() string {
	return fmt.Sprintf("there is not avaiable worker with indigrient %s", we.name)
}
