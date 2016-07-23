package manager

import (
	"regexp"
	"strings"
	"os/exec"
	"bufio"
	"io"
)

type Task interface {
	Inputs() []Pair
	InputOf(name string) Pair
	Outputs() []Pair
	OutputOf(name string) Pair
	Run(runner Runner)
}

type Pair interface {
	Key() string
	Value() string
	Set(val string)
}

type Runner interface {
	Run(path string, args []string)
	Stdout() (chan string)
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

func (t *TaskImpl) Run(runner Runner) {
	for k, v := range t.inputMap {
		if v == "" {
			panic("Parameter " + k + " not fulfilled")
		}
	}
	for i, arg := range t.args {
		if strings.HasPrefix(arg, "$") {
			varname := strings.TrimPrefix(arg, "$")
			t.args[i] = t.inputMap[varname]
		}
	}
	runner.Run(t.path, t.args)
	for output := range runner.Stdout() {
		for key, regex := range t.filters {
			if match := regex.FindString(output) ; match != "" {
				t.outputMap[key] = match
			}
		}
	}
}

func NewTask(template string, inputs []string, outputs map[string]string) *TaskImpl {
	splitted := strings.Split(template, " ")
	path := splitted[0]
	args := splitted[1:]
	filters := make(map[string]*regexp.Regexp, len(outputs))
	outputMap := make(map[string]string, len(outputs))
	for k, v := range outputs {
		reg:= regexp.MustCompile(v)
		filters[k] = reg
		outputMap[k] = ""
	}
	inputMap := make(map[string]string, len(inputs))
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

type RunnerImpl struct {
	cmd *exec.Cmd
	output chan string
}

func (r *RunnerImpl) Run(path string, args []string) {
	r.cmd = exec.Command(path, args...)
	reader, err := r.cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	err = r.cmd.Start()
	if err != nil {
		panic(err)
	}
	r.output = make(chan string)
	go r.redirectOuput(reader)
} 

func (r *RunnerImpl) redirectOuput(reader io.Reader) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		r.output <- scanner.Text()
	}
	close(r.output)
}

func (r *RunnerImpl) Stdout() (chan string) {
	return r.output
}

func NewRunner() Runner {
	return &RunnerImpl{}
}