package spirit

import "time"

type Cocktail struct{
	Name string `json:"name"`
	Ingredient []Ingredient `json:"ingredient"`
}

type Ingredient struct {
	Type string `json:"type"`
	Duration time.Duration `json:"duration"`
	Order int `json:"order"`
}




