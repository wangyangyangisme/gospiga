// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package gospiga

import (
	"fmt"
	"io"
	"strconv"
)

type Ingredient struct {
	UID      string `json:"uid"`
	Name     string `json:"name"`
	Quantity *int   `json:"quantity"`
}

type NewIngredient struct {
	Name     string `json:"name"`
	Quantity *int   `json:"quantity"`
}

type NewRecipe struct {
	Title       string           `json:"title"`
	Subtitle    string           `json:"subtitle"`
	Description *string          `json:"description"`
	Ingredient  []*NewIngredient `json:"ingredient"`
	Step        []*NewStep       `json:"step"`
	Conclusion  *string          `json:"conclusion"`
}

type NewStep struct {
	Index   int    `json:"index"`
	Excerpt string `json:"excerpt"`
	Text    string `json:"text"`
}

type Step struct {
	UID     string `json:"uid"`
	Index   int    `json:"index"`
	Excerpt string `json:"excerpt"`
	Text    string `json:"text"`
}

type UpIngredient struct {
	UID      *string `json:"uid"`
	Name     string  `json:"name"`
	Quantity *int    `json:"quantity"`
}

type UpRecipe struct {
	UID         string          `json:"uid"`
	Title       string          `json:"title"`
	Subtitle    string          `json:"subtitle"`
	Description *string         `json:"description"`
	Ingredient  []*UpIngredient `json:"ingredient"`
	Step        []*UpStep       `json:"step"`
	Conclusion  *string         `json:"conclusion"`
}

type UpStep struct {
	UID     *string `json:"uid"`
	Index   int     `json:"index"`
	Excerpt string  `json:"excerpt"`
	Text    string  `json:"text"`
}

type Unit string

const (
	UnitGrams       Unit = "GRAMS"
	UnitKilograms   Unit = "KILOGRAMS"
	UnitOunces      Unit = "OUNCES"
	UnitPounds      Unit = "POUNDS"
	UnitLiters      Unit = "LITERS"
	UnitMilliliters Unit = "MILLILITERS"
)

var AllUnit = []Unit{
	UnitGrams,
	UnitKilograms,
	UnitOunces,
	UnitPounds,
	UnitLiters,
	UnitMilliliters,
}

func (e Unit) IsValid() bool {
	switch e {
	case UnitGrams, UnitKilograms, UnitOunces, UnitPounds, UnitLiters, UnitMilliliters:
		return true
	}
	return false
}

func (e Unit) String() string {
	return string(e)
}

func (e *Unit) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Unit(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Unit", str)
	}
	return nil
}

func (e Unit) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
