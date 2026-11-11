package converter

import (
	"errors"
	"math"
	"testing"
)

func TestConvert(t *testing.T) {
	tests := []struct {
		name    string
		value   float64
		from    string
		to      string
		want    float64
		wantErr error
	}{
		{name: "1 mètre en centimètres", value: 1, from: "mètre", to: "centimètre", want: 100, wantErr: nil},
		{name: "1 mètre en pouces", value: 1, from: "mètre", to: "pouce", want: 39.37007874015748, wantErr: nil},
		{name: "mètre vers kilogramme (incompatible)", value: 1, from: "mètre", to: "kilogramme", want: 0, wantErr: ErrIncompatibleUnits},
		{name: "vers unité inconnue", value: 1, from: "mètre", to: "unité_inconnue", want: 0, wantErr: ErrUnknownUnit},
		{name: "depuis unité inconnue", value: 1, from: "unité_inconnue", to: "mètre", want: 0, wantErr: ErrUnknownUnit},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Convert(tc.value, tc.from, tc.to)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("Convert() error = %v, wantErr %v", err, tc.wantErr)
			}
			if tc.wantErr != nil {
				return
			}
			if math.Abs(got-tc.want) > 1e-9 {
				t.Errorf("Convert() = %v, want %v", got, tc.want)
			}
		})
	}
}
