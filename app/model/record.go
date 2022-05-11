package model

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

var periodTemplate = regexp.MustCompile("(?:(Monday|Tuesday|Wednesday|Thursday|Friday|Saturday|Sunday)) ([1-9]|1[0-2])AM - ([1-9]|1[0-2])PM")

type Recipe struct {
	fullName string
	names    []string
}

type DeliveryPeriod struct {
	weekday string
	Range   TimeRange
}

type TimeRange struct {
	AM uint
	PM uint
}

type Record struct {
	Postcode string         `json:"postcode"`
	Recipe   Recipe         `json:"recipe"`
	Delivery DeliveryPeriod `json:"delivery"`
}

func NewRecipe(name string) Recipe {
	fullName := strings.Trim(name, "\"")
	return Recipe{
		fullName: fullName,
		names:    strings.Split(fullName, " "),
	}
}

func (r Recipe) String() string {
	return r.fullName
}

func (r Recipe) ContainsOneOf(words []string) bool {
	for _, word := range words {
		for _, name := range r.names {
			if strings.ToLower(name) == strings.ToLower(word) {
				return true
			}
		}
	}
	return false
}

func (r *Recipe) UnmarshalJSON(b []byte) error {
	*r = NewRecipe(string(b))
	return nil
}

func (t TimeRange) In(period TimeRange) bool {
	return t.AM >= period.AM && t.PM <= period.PM
}

func NewDeliveryPeriod(data string) (DeliveryPeriod, error) {
	p := DeliveryPeriod{}
	if !periodTemplate.MatchString(data) {
		return p, errors.New("bad delivery format")
	}
	subMatch := periodTemplate.FindStringSubmatch(data)
	p.weekday = subMatch[1]

	am, err := strconv.Atoi(subMatch[2])
	if err != nil {
		return p, err
	}
	p.Range.AM = uint(am)

	pm, err := strconv.Atoi(subMatch[3])
	if err != nil {
		return p, err
	}
	p.Range.PM = uint(pm)
	return p, nil
}

func (p *DeliveryPeriod) UnmarshalJSON(b []byte) error {
	var err error
	*p, err = NewDeliveryPeriod(string(b))
	return err
}
