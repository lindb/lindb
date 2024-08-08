package plan

import (
	"fmt"
	"testing"
)

func TestAssignments(t *testing.T) {
	var assignments Assignments
	assignments = assignments.Add([]*Symbol{{}})
	fmt.Println(assignments)
}
