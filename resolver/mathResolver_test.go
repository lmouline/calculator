package resolver

import (
	"testing"
	"go/constant"
	"errors"
)

type ExpectedValues struct {
	value constant.Value
	err	error
}

var RESULTS = map[string]ExpectedValues{
	"":                               {nil, errors.New("")},
	"func tt(){}":                    {value: nil, err: errors.New("")},
	"false":                          {value: constant.MakeBool(false), err: nil},
	"TRUE":                           {value: constant.MakeBool(true), err: nil},
	"true":                           {value: constant.MakeBool(true), err: nil},
	"FALSE":                          {value: constant.MakeBool(false), err: nil},
	"((1 + 1) - (4 * 5))*0 + 2 == 5": {value: constant.MakeBool(false), err: nil},
	"((1 + 1) - (4 * 5))*0 + 2":      {value: constant.MakeInt64(2), err: nil},
	"(1+1.1)*2":                      {value: constant.MakeFloat64(4.2), err: nil},
	//"3pi^2":                      {value: constant.MakeFloat64(29.608813203268074), err: nil},
}

func TestResolveBoolExpr(t *testing.T)  {
	for toEvaluate, expected := range RESULTS  {
		res, err := Resolve(toEvaluate)
		if (err == nil && expected.err != nil) || (err != nil && expected.err == nil){
			t.Errorf("For expression '%s': Error mismatch",toEvaluate)
		}

		if (res == nil && expected.value != nil) || (err != nil && expected.err == nil){
			t.Errorf("For expression '%s': [Expected] %t != %t [Actual]", toEvaluate, expected.value, res)
		}

		if res != nil && (res.String() != expected.value.String()) {
			t.Errorf("For expression '%s': [Expected] %t != %t [Actual]", toEvaluate, expected.value, res)
		}
	}
}
