package hw09structvalidator

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type (
	IntCheck struct {
		MinInt int `validate:"min:500"`
		MaxInt int `validate:"max:100"`
		InInt  int `validate:"in:10,100"`
	}

	StringCheck struct {
		SmallStr string `validate:"max:10"`
		RegexStr string `validate:"regexp:^[abc]+$"`
		InStr    string `validate:"in:testIn1,testIn2"`
		LenStr   string `validate:"len:10"`
	}

	SliceCheck struct {
		SliceStrIn []string `validate:"in:first,second,third"`
		SliceIntIn []int    `validate:"in:1,10"`
	}

	Complex struct {
		IntCheck        IntCheck    `validate:"nested"`
		StringCheck     StringCheck `validate:"nested"`
		Nested          Nested      `validate:"nested"`
		SliceCheck      SliceCheck  `validate:"nested"`
		WithJSON        int         `json:"test" validate:"in:1,10"`
		WithJSON2       int         `json:"jsonjson,omitempty" validate:"in:1,10"`
		ignoreThisField string      `validate:"len:10"`
	}

	Nested struct {
		NestedStr    string       `validate:"min:10|max:15"`
		NestedInt    int          `validate:"max:15"`
		NestedNested NestedNested `validate:"nested"`
	}

	NestedNested struct {
		NestedNestedStr string `validate:"min:20"`
	}
)

func TestValidate(t *testing.T) {
	passInts := IntCheck{
		MinInt: 600,
		MaxInt: 90,
		InInt:  50,
	}
	passStrings := StringCheck{
		SmallStr: "sfasaffas",
		RegexStr: "abcbacbcbacac",
		InStr:    "testIn1",
		LenStr:   "10lengthex",
	}
	passSlice := SliceCheck{
		SliceStrIn: []string{"first", "second", "third"},
		SliceIntIn: []int{2, 4, 5, 6, 1, 6, 5},
	}
	nestedNested := NestedNested{
		NestedNestedStr: "12345678901234567890",
	}
	nested := Nested{
		NestedStr:    "minimum10sym",
		NestedInt:    14,
		NestedNested: nestedNested,
	}
	complexPass := Complex{
		IntCheck:        passInts,
		StringCheck:     passStrings,
		SliceCheck:      passSlice,
		Nested:          nested,
		WithJSON:        9,
		WithJSON2:       2,
		ignoreThisField: "ignoreignoreignoreignoreignore",
	}

	wrongInts := IntCheck{
		MinInt: 2,
		MaxInt: 990,
		InInt:  1,
	}
	wrongIntsErrors := ValidationErrors{
		ValidationError{
			Field: "MinInt",
			Err:   fmt.Errorf(minIntWrongMes, 500),
		},
		ValidationError{
			Field: "MaxInt",
			Err:   fmt.Errorf(maxIntWrongMes, 100),
		},
		ValidationError{
			Field: "InInt",
			Err:   fmt.Errorf(inWrongMes, "10,100"),
		},
	}

	wrongStrings := StringCheck{
		SmallStr: "veribigstringveribigstringveribigstring",
		RegexStr: "wrongregex",
		InStr:    "notInStr",
		LenStr:   "not10lengthhh",
	}
	wrongStringsErrors := ValidationErrors{
		ValidationError{
			Field: "SmallStr",
			Err:   fmt.Errorf(maxWrongMes, 10),
		},
		ValidationError{
			Field: "RegexStr",
			Err:   fmt.Errorf(regexpWrongMes),
		},
		ValidationError{
			Field: "InStr",
			Err:   fmt.Errorf(inWrongMes, "testIn1,testIn2"),
		},
		ValidationError{
			Field: "LenStr",
			Err:   fmt.Errorf(lenWrongMes, 10),
		},
	}

	wrongSlice := SliceCheck{
		SliceStrIn: []string{"NotFirst", "NotSecond", "NotThird"},
		SliceIntIn: []int{12, 24, 65, 56, 11},
	}
	wrongSliceErrors := ValidationErrors{
		ValidationError{
			Field: "SliceStrIn",
			Err:   fmt.Errorf(inWrongMes, "first,second,third"),
		},
		ValidationError{
			Field: "SliceIntIn",
			Err:   fmt.Errorf(inWrongMes, "1,10"),
		},
	}

	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in:          complexPass,
			expectedErr: nil,
		},
		{
			in:          wrongInts,
			expectedErr: wrongIntsErrors,
		},
		{
			in:          wrongStrings,
			expectedErr: wrongStringsErrors,
		},
		{
			in:          wrongSlice,
			expectedErr: wrongSliceErrors,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)

			require.Equal(t, tt.expectedErr, err)
		})
	}
}
