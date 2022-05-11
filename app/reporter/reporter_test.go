package reporter

import (
	"github.com/rs/zerolog"
	"github.com/sanches1984/json-parser/app/config"
	"github.com/sanches1984/json-parser/app/model"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestReporter_Process(t *testing.T) {
	reporter := New(config.Config{
		Filter: config.SearchFilter{
			Postcode:    "10127",
			AM:          10,
			PM:          5,
			RecipeNames: []string{"Steak"},
		},
	}, zerolog.New(os.Stderr))

	reporter.processRecord(model.Record{
		Postcode: "10125",
		Recipe:   model.NewRecipe("Creamy Dill Chicken"),
		Delivery: newDeliveryPeriod("Monday 9AM - 8PM"),
	})
	reporter.processRecord(model.Record{
		Postcode: "10127",
		Recipe:   model.NewRecipe("Speedy Steak Fajitas"),
		Delivery: newDeliveryPeriod("Monday 9AM - 8PM"),
	})
	reporter.processRecord(model.Record{
		Postcode: "10127",
		Recipe:   model.NewRecipe("Cherry Balsamic Pork Chops"),
		Delivery: newDeliveryPeriod("Monday 11AM - 4PM"),
	})
	reporter.processRecord(model.Record{
		Postcode: "10129",
		Recipe:   model.NewRecipe("Cherry Balsamic Pork Chops"),
		Delivery: newDeliveryPeriod("Monday 9AM - 4PM"),
	})

	require.Equal(t, map[string]struct{}{"Speedy Steak Fajitas": {}}, reporter.foundRecipes)
	require.Equal(t, map[string]uint{"Speedy Steak Fajitas": 1, "Creamy Dill Chicken": 1, "Cherry Balsamic Pork Chops": 2}, reporter.uniqueRecipe)
	require.Equal(t, model.PostcodeCount{Postcode: "10127", Count: 1}, reporter.searchPostcode)
	require.Equal(t, model.PostcodeCount{Postcode: "10127", Count: 2}, reporter.busiestPostcode)
	require.Equal(t, map[string]uint{"10125": 1, "10127": 2, "10129": 1}, reporter.postcodeDeliveries)

}

func TestReporter_MakeReport(t *testing.T) {
	reporter := &Reporter{
		logger: zerolog.New(os.Stderr),
		filter: filter{
			Postcode:       "10127",
			DeliveryPeriod: model.TimeRange{AM: 12, PM: 4},
			RecipeNames:    []string{"Steak"},
		},
		uniqueRecipe: map[string]uint{
			"Creamy Dill Chicken":        12,
			"Speedy Steak Fajitas":       23,
			"Cherry Balsamic Pork Chops": 41,
			"Hot Steak":                  17,
			"Some Unique Recipe":         1,
		},
		postcodeDeliveries: map[string]uint{
			"10120": 11,
			"11232": 7,
			"10134": 5,
		},
		searchPostcode: model.PostcodeCount{
			Postcode: "10127",
			Count:    4,
		},
		busiestPostcode: model.PostcodeCount{
			Postcode: "10120",
			Count:    11,
		},
		foundRecipes: map[string]struct{}{
			"Speedy Steak Fajitas": {},
			"Hot Steak":            {},
		},
	}

	expected := model.Report{
		UniqueRecipeCount: 1,
		CountPerRecipe: []model.RecipeCount{
			{
				Recipe: "Cherry Balsamic Pork Chops",
				Count:  41,
			},
			{
				Recipe: "Creamy Dill Chicken",
				Count:  12,
			},
			{
				Recipe: "Hot Steak",
				Count:  17,
			},
			{
				Recipe: "Some Unique Recipe",
				Count:  1,
			},
			{
				Recipe: "Speedy Steak Fajitas",
				Count:  23,
			},
		},
		BusiestCode: model.PostcodeCount{
			Postcode: "10120",
			Count:    11,
		},
		CountPerPostcodeAndTime: model.PostcodeDeliveries{
			Postcode:      "10127",
			From:          "12AM",
			To:            "4PM",
			DeliveryCount: 4,
		},
		MatchByName: []string{"Hot Steak", "Speedy Steak Fajitas"},
	}

	require.Equalf(t, expected, reporter.MakeReport(), "wrong report")
}

func newDeliveryPeriod(data string) model.DeliveryPeriod {
	p, _ := model.NewDeliveryPeriod(data)
	return p
}
