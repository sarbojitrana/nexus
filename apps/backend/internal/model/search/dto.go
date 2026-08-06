package search

import "github.com/go-playground/validator/v10"

type SearchQuery struct {
	Query string `query:"q" validate:"required,min=1,max=200"`
}

func (p *SearchQuery) Validate() error {
	return validator.New().Struct(p)
}
