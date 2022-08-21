package core

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/dangduoc08/go-go/helper"
)

const VAR_SYMBOL = "$"

type Var[T interface{}] struct{ KeyValue (map[string]T) }

func NewVar(s string) (string, *Var[interface{}]) {
	matchVarPatternReg := regexp.MustCompile(fmt.Sprintf(`\%v(.*?)\%v`, helper.OPEN_CURLY_BRACKET, helper.CLOSE_CURLY_BRACKET))
	v := &Var[interface{}]{
		KeyValue: make(map[string]interface{}),
	}

	for _, valB := range matchVarPatternReg.FindAll([]byte(s), -1) {
		valStr := string(valB)
		s = strings.Replace(s, valStr, VAR_SYMBOL, 1)
		varKey := helper.RemoveAtEnd(helper.RemoveAtBegin(valStr, helper.OPEN_CURLY_BRACKET), helper.CLOSE_CURLY_BRACKET)
		v.Set(varKey, nil)
	}

	return s, v
}

func (v *Var[T]) Get(k string) T {
	return v.KeyValue[k]
}

func (v *Var[T]) Set(k string, val T) {
	v.KeyValue[k] = val
}

func (v *Var[T]) Delete(k string) {
	delete(v.KeyValue, k)
}
