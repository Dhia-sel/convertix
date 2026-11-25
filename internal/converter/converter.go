// Package converter converts values between units of length and mass.
package converter

import (
	"errors"
	"fmt"
)

var (
	// ErrUnknownUnit is returned when a unit is not recognised.
	ErrUnknownUnit = errors.New("l'unité recherchée est inconnue")
	// ErrIncompatibleUnits is returned when the two units are of different kinds.
	ErrIncompatibleUnits = errors.New("les unités sont incompatibles")
)

type unitInfo struct {
	factor   float64
	category string
}

var units = map[string]unitInfo{
	"millimètre":      {factor: 0.001, category: "length"},
	"centimètre":      {factor: 0.01, category: "length"},
	"kilomètre":       {factor: 1000.0, category: "length"},
	"mètre":           {factor: 1.0, category: "length"},
	"pouce":           {factor: 0.0254, category: "length"},
	"pied":            {factor: 0.3048, category: "length"},
	"yard":            {factor: 0.9144, category: "length"},
	"rod":             {factor: 5.0292, category: "length"},
	"furlong":         {factor: 201.168, category: "length"},
	"mille terrestre": {factor: 1609.344, category: "length"},
	"mille marin":     {factor: 1852.0, category: "length"},
	"lieue":           {factor: 4828.032, category: "length"},

	"gramme":       {factor: 0.001, category: "mass"},
	"tonne":        {factor: 1000.0, category: "mass"},
	"kilogramme":   {factor: 1.0, category: "mass"},
	"grain":        {factor: 0.00006479891, category: "mass"},
	"drachme":      {factor: 0.00177184519531, category: "mass"},
	"once":         {factor: 0.028349523125, category: "mass"},
	"livre":        {factor: 0.45359237, category: "mass"},
	"stone":        {factor: 6.35029318, category: "mass"},
	"quintal":      {factor: 100.0, category: "mass"},
	"tonne courte": {factor: 907.18474, category: "mass"},
	"tonne longue": {factor: 1016.04691, category: "mass"},
}

// Convert converts value from the unit "from" to the unit "to".
// It returns ErrUnknownUnit if a unit is not recognised, and
// ErrIncompatibleUnits if the two units are of different kinds.
func Convert(value float64, from, to string) (float64, error) {
	unitFrom, okFrom := units[from]
	if !okFrom {
		return 0, fmt.Errorf("%w: %q", ErrUnknownUnit, from)
	}

	unitTo, okTo := units[to]
	if !okTo {
		return 0, fmt.Errorf("%w: %q", ErrUnknownUnit, to)
	}

	if unitFrom.category != unitTo.category {
		return 0, fmt.Errorf("%w: %s (%s) et %s (%s)", ErrIncompatibleUnits, from, unitFrom.category, to, unitTo.category)
	}

	return value * unitFrom.factor / unitTo.factor, nil
}
