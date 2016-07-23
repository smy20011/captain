package manager

import (
	"testing"
	"reflect"
)

func assert(t *testing.T, pred bool, desp string) {
	if !pred {
		t.Log("Test failed: " + desp)
		t.Fail()
	}
}

func assertEquals(t *testing.T, excepted, actuall interface{}, desp string) {
	assert(t, reflect.DeepEqual(excepted, actuall), desp)
}

func SimpleTask() Task {
	return NewTask("vim simple.txt $input", []string{"input"}, map[string]string{"progress": "\\d+"})
}

func TestCreateTask(t *testing.T) {
	task := SimpleTask()
	assertEquals(t, len(task.Inputs()), 1, "Should have one input")
	assertEquals(t, len(task.Outputs()), 1, "Should have zero output")
}

type DummyRunner struct {
	output chan string
	path string
	args []string
}

func (r *DummyRunner) Run(path string, args []string) {
	r.path = path
	r.args = args
}

func (r *DummyRunner) Stdout() (chan string) {
	return r.output
}

func NewDummyRunner() *DummyRunner {
	return &DummyRunner{make(chan string), "", []string{}}
}

func TestRunTask(t *testing.T) {
	task := SimpleTask()
	runner := NewDummyRunner()
	task.InputOf("input").Set("something")
	go task.Run(runner)
	assertEquals(t, task.OutputOf("progress").Value(), "", "No Output")
	runner.Stdout() <- "123"
	assertEquals(t, task.OutputOf("progress").Value(), "123", "Capture output correctly")
	runner.Stdout() <- "abc"
	assertEquals(t, task.OutputOf("progress").Value(), "123", "Ignore unrelated output")
	close(runner.Stdout())
	assertEquals(t, runner.path, "vim", "Parse path correctly")
	assertEquals(t, runner.args, []string{"simple.txt", "something"}, "Replce args correctly")
}

func TestRunnerImpl(t *testing.T) {
	runner := NewRunner()
	runner.Run("echo", []string{"Hello"})
	output := <- runner.Stdout()
	assertEquals(t, "Hello", output, "Should say hello")
}