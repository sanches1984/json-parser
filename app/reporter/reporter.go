package reporter

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/sanches1984/json-parser/app/config"
	"github.com/sanches1984/json-parser/app/model"
	"io"
	"sort"
)

const logCount = 100000

type Reporter struct {
	filter filter
	logger zerolog.Logger

	uniqueRecipe       map[string]uint
	postcodeDeliveries map[string]uint
	searchPostcode     model.PostcodeCount
	busiestPostcode    model.PostcodeCount
	foundRecipes       map[string]struct{}
}

type filter struct {
	Postcode       string
	DeliveryPeriod model.TimeRange
	RecipeNames    []string
}

func New(config config.Config, logger zerolog.Logger) *Reporter {
	return &Reporter{
		filter: filter{
			Postcode: config.Filter.Postcode,
			DeliveryPeriod: model.TimeRange{
				AM: config.Filter.AM,
				PM: config.Filter.PM,
			},
			RecipeNames: config.Filter.RecipeNames,
		},
		logger:             logger,
		uniqueRecipe:       make(map[string]uint),
		postcodeDeliveries: make(map[string]uint),
		searchPostcode:     model.PostcodeCount{Postcode: config.Filter.Postcode},
		busiestPostcode:    model.PostcodeCount{},
		foundRecipes:       make(map[string]struct{}),
	}
}

func (r *Reporter) Read(file io.Reader) error {
	dec := json.NewDecoder(file)

	if _, err := dec.Token(); err != nil {
		return err
	}

	count := 0
	for dec.More() {
		var record model.Record
		err := dec.Decode(&record)
		if err != nil {
			return err
		}

		r.processRecord(record)
		count++

		if count%logCount == 0 {
			r.logger.Debug().Int("count", count).Msg("parsing")
		}
	}

	if _, err := dec.Token(); err != nil {
		return err
	}

	r.logger.Debug().Int("count", count).Msg("parsed")
	return nil
}

func (r Reporter) MakeReport() model.Report {
	return model.Report{
		UniqueRecipeCount: r.getUniqueCount(),
		CountPerRecipe:    r.getCountPerRecipe(),
		BusiestCode:       r.busiestPostcode,
		MatchByName:       r.getFoundRecipes(),
		CountPerPostcodeAndTime: model.PostcodeDeliveries{
			Postcode:      r.searchPostcode.Postcode,
			From:          fmt.Sprintf("%dAM", r.filter.DeliveryPeriod.AM),
			To:            fmt.Sprintf("%dPM", r.filter.DeliveryPeriod.PM),
			DeliveryCount: r.searchPostcode.Count,
		},
	}
}

func (r *Reporter) processRecord(record model.Record) {
	// count unique recipes' occurrences
	r.uniqueRecipe[record.Recipe.String()]++

	// count deliveries on a searching postcode within period
	if record.Postcode == r.filter.Postcode && record.Delivery.Range.In(r.filter.DeliveryPeriod) {
		r.searchPostcode.Count++
	}

	// find busiest postcode
	r.postcodeDeliveries[record.Postcode]++
	if r.busiestPostcode.Count == 0 {
		r.busiestPostcode.Postcode = record.Postcode
		r.busiestPostcode.Count = 1
	} else if r.busiestPostcode.Count < r.postcodeDeliveries[record.Postcode] {
		r.busiestPostcode.Postcode = record.Postcode
		r.busiestPostcode.Count = r.postcodeDeliveries[record.Postcode]
	}

	// find recipes by names
	if record.Recipe.ContainsOneOf(r.filter.RecipeNames) {
		if _, ok := r.foundRecipes[record.Recipe.String()]; !ok {
			r.foundRecipes[record.Recipe.String()] = struct{}{}
		}
	}
}

func (r Reporter) getUniqueCount() uint {
	var count uint
	for _, entries := range r.uniqueRecipe {
		if entries == 1 {
			count++
		}
	}
	return count
}

func (r Reporter) getCountPerRecipe() []model.RecipeCount {
	recipes := make([]model.RecipeCount, 0, len(r.uniqueRecipe))
	for recipe, count := range r.uniqueRecipe {
		recipes = append(recipes, model.RecipeCount{
			Recipe: recipe,
			Count:  count,
		})
	}

	sort.Slice(recipes, func(i, j int) bool {
		return recipes[i].Recipe < recipes[j].Recipe
	})
	return recipes
}

func (r Reporter) getFoundRecipes() []string {
	recipes := make([]string, 0, len(r.foundRecipes))
	for recipe := range r.foundRecipes {
		recipes = append(recipes, recipe)
	}

	sort.Strings(recipes)
	return recipes
}
