package model

import "encoding/json"

type Report struct {
	UniqueRecipeCount       uint               `json:"unique_recipe_count"`
	CountPerRecipe          []RecipeCount      `json:"count_per_recipe"`
	BusiestCode             PostcodeCount      `json:"busiest_code"`
	CountPerPostcodeAndTime PostcodeDeliveries `json:"count_per_postcode_and_time"`
	MatchByName             []string           `json:"match_by_name"`
}

type RecipeCount struct {
	Recipe string `json:"recipe"`
	Count  uint   `json:"count"`
}

type PostcodeCount struct {
	Postcode string `json:"postcode"`
	Count    uint   `json:"delivery_count"`
}

type PostcodeDeliveries struct {
	Postcode      string `json:"postcode"`
	From          string `json:"from"`
	To            string `json:"to"`
	DeliveryCount uint   `json:"delivery_count"`
}

func (o Report) ToJSON() ([]byte, error) {
	return json.Marshal(&o)
}
