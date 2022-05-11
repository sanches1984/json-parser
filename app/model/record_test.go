package model

import (
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTimeRange_In(t *testing.T) {
	mainRange := TimeRange{AM: 9, PM: 3}
	cases := []struct {
		timeRange TimeRange
		hit       bool
	}{
		{
			timeRange: TimeRange{AM: 9, PM: 2},
			hit:       true,
		},
		{
			timeRange: TimeRange{AM: 11, PM: 1},
			hit:       true,
		},
		{
			timeRange: TimeRange{AM: 12, PM: 4},
			hit:       false,
		},
		{
			timeRange: TimeRange{AM: 12, PM: 4},
			hit:       false,
		},
		{
			timeRange: TimeRange{AM: 8, PM: 2},
			hit:       false,
		},
		{
			timeRange: TimeRange{AM: 8, PM: 5},
			hit:       false,
		},
	}

	for n, c := range cases {
		require.Equalf(t, c.hit, c.timeRange.In(mainRange), "wrong result in case %d", n)
	}
}

func TestRecipe_ContainsOneOf(t *testing.T) {
	names := []string{"Potato", "Veggie", "Mushroom"}
	cases := []struct {
		recipe Recipe
		hit    bool
	}{
		{
			recipe: Recipe{
				names: []string{"Some", "Wrong", "Words", "Example"},
			},
			hit: false,
		},
		{
			recipe: Recipe{
				names: []string{"Some", "Veggie", "Words", "Example"},
			},
			hit: true,
		},
		{
			recipe: Recipe{
				names: []string{"Some", "Potato"},
			},
			hit: true,
		},
		{
			recipe: Recipe{
				names: []string{"Mushroom", "Potato", "Words", "Veggie"},
			},
			hit: true,
		},
		{
			recipe: Recipe{
				names: []string{},
			},
			hit: false,
		},
	}

	for n, c := range cases {
		require.Equalf(t, c.hit, c.recipe.ContainsOneOf(names), "wrong result in case %d", n)
	}
}

func TestRecipe_UnmarshalJSON(t *testing.T) {
	cases := []struct {
		data     []byte
		expected Recipe
	}{
		{
			data: []byte("Mushroom Potato Words Veggie"),
			expected: Recipe{
				fullName: "Mushroom Potato Words Veggie",
				names:    []string{"Mushroom", "Potato", "Words", "Veggie"},
			},
		},
		{
			data: []byte("OneWord"),
			expected: Recipe{
				fullName: "OneWord",
				names:    []string{"OneWord"},
			},
		},
		{
			data: []byte(""),
			expected: Recipe{
				fullName: "",
				names:    []string{""},
			},
		},
	}

	for n, c := range cases {
		var r Recipe
		require.NoErrorf(t, r.UnmarshalJSON(c.data), "unexpected error in case %d", n)
		require.Equalf(t, c.expected, r, "wrong result in case %d", n)
	}
}

func TestDeliveryPeriod_UnmarshalJSON_Success(t *testing.T) {
	cases := []struct {
		data     []byte
		expected DeliveryPeriod
	}{
		{
			data: []byte("Wednesday 9AM - 10PM"),
			expected: DeliveryPeriod{
				weekday: "Wednesday",
				Range:   TimeRange{AM: 9, PM: 10},
			},
		},
		{
			data: []byte("Monday 12AM - 5PM"),
			expected: DeliveryPeriod{
				weekday: "Monday",
				Range:   TimeRange{AM: 12, PM: 5},
			},
		},
	}

	for n, c := range cases {
		var p DeliveryPeriod
		require.NoErrorf(t, p.UnmarshalJSON(c.data), "unexpected error in case %d", n)
		require.Equalf(t, c.expected, p, "wrong result in case %d", n)
	}
}

func TestDeliveryPeriod_UnmarshalJSON_Error(t *testing.T) {
	err := errors.New("bad delivery format")
	cases := []struct {
		data []byte
	}{
		{
			data: []byte("Wedney 9AM - 10PM"),
		},
		{
			data: []byte("Wednesday 9AM-10PM"),
		},
		{
			data: []byte("Monday 14AM - 15PM"),
		},
	}

	for n, c := range cases {
		var p DeliveryPeriod
		require.EqualErrorf(t, p.UnmarshalJSON(c.data), err.Error(), "unexpected error in case %d", n)
	}
}
