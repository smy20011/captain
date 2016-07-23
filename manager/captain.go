package manager

import (
	"regexp"
	"strings"
)

type Task interface {
	Inputs() []Pair
	InputOf(name string) Pair
	Outputs() Pair
	OutputOf(name string) Pair
	Run()
}

type Pair interface {
	Key() string
	Value() string
	Set(val string)
}

type MapPair struct {
	m   map[string]string
	key string
}

func (m *MapPair) Key() string {
	return m.key
}

func (m *MapPair) Value() string {
	return m.m[m.key]
}

func (m *MapPair) Set(val string) {
	m.m[m.key] = val
}

func NewMapPair(mapObj map[string]string, key string) *MapPair {
	return &MapPair{mapObj, key}
}

func GetPairs(mapObj map[string]string) []Pair {
	result := make([]Pair, len(mapObj))
	index := 0
	for k, _ := range mapObj {
		result[index] = NewMapPair(mapObj, k)
		index++
	}
	return result
}

type TaskImpl struct {
	inputMap  map[string]string
	outputMap map[string]string
	path      string
	args      []string
	filters   map[string]*regexp.Regexp
}

func (t *TaskImpl) Inputs() []Pair {
	return GetPairs(t.inputMap)
}

func (t *TaskImpl) InputOf(name string) Pair {
	return NewMapPair(t.inputMap, name)
}

func (t *TaskImpl) Outputs() []Pair {
	return GetPairs(t.outputMap)
}

func (t *TaskImpl) OutputOf(name string) Pair {
	return NewMapPair(t.outputMap, name)
}

func (t *TaskImpl) Run() {
	panic("not implemented")
}

func NewTask(template string, inputs []string, outputs map[string]string) *TaskImpl {
	splitted := strings.Split(template, " ")
	path := splitted[0]
	args := splitted[1:]
	filters := make(map[string]*regexp.Regexp)
	outputMap := make(map[string]string)
	for k, v := range outputs {
		if reg, err := regexp.Compile(v); err == nil {
			filters[k] = reg
			outputMap[k] = ""
		} else {
			panic(err)
		}
	}
	inputMap := make(map[string]string)
	for _, k := range inputs {
		inputMap[k] = ""
	}
	return &TaskImpl{
		inputMap:  inputMap,
		outputMap: outputMap,
		path:      path,
		args:      args,
		filters:   filters,
	}
}
