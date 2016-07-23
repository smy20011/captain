package manager

import (
	"testing"
)

func assert(t *testing.T, pred bool, desp string) {
	if !pred {
		t.Log("Test failed: " + desp)
		t.Fail()
	}
}

func TestCreateTask(t *testing.T) {
	task := NewTask("vim simple.txt $input", []string{"input"}, map[string]string{})
	assert(t, len(task.Inputs()) == 1, "Should have one input")
	assert(t, len(task.Outputs()) == 0, "Should have zero output")
	assert(t, len(task.filters) == 0, "Should have zero filters")
}
