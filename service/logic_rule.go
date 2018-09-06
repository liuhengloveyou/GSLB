package service

import (
	"fmt"
)

type LogicRule interface {
	Resolv(k, v string) error
}

type LogicRuleType func() (LogicRule, error)

var logicRules map[string]LogicRuleType = make(map[string]LogicRuleType)

func RegisterLogicRule(name string, newFunc LogicRuleType) {
	if newFunc == nil {
		panic("Register LogicRule nil.")
	}

	if _, ok := logicRules[name]; ok {
		panic("Register LogicRule duplicate for " + name)
	}

	logicRules[name] = newFunc
}

func ResolvLogicRule(k, v string) (err error) {
	newFunc, ok := logicRules[k]
	if !ok {
		return fmt.Errorf("No Logic Rule types " + k)
	}

	logic, err := newFunc()
	if !ok {
		return fmt.Errorf("New Logic Rule ERR " + err.Error())
	}

	return logic.Resolv(k, v)
}
