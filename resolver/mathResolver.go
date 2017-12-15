package resolver

import (
	"go/constant"
	"go/token"
	"go/types"
	"strings"
)

func Resolve(expression string) (val constant.Value, err error) {
	exp := strings.ToLower(expression)
	fset := token.NewFileSet()
	typeAndVal, err := types.Eval(fset,nil,0,exp)
	val = typeAndVal.Value

	return
}



