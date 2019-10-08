package tree

import (
	"bytes"
	"fmt"
	"io"

	pomo "github.com/kevinschoon/pomo/pkg"
	"github.com/kevinschoon/pomo/pkg/internal/format"
)

const (
	emptySpace   = "    "
	middleItem   = "├──"
	continueItem = "│   "
	lastItem     = "└──"
)

// Tree prints a hierarchy of tasks
// in a tree.
type Tree struct {
	pomo.Task
	ShowPomodoros bool
}

func (t Tree) next(value bool, depth []bool) (result []bool) {
	for _, value := range depth {
		result = append(result, value)
	}
	result = append(result, value)
	return result
}

func (t Tree) fill(w io.Writer, depth []bool) {
	for i := 0; i < len(depth); i++ {
		if depth[i] {
			fmt.Fprintf(w, continueItem)
		} else {
			fmt.Fprintf(w, emptySpace)
		}
	}
}

// Write writes the Tree representation of a Task hierarchy
// to the io.Writer
func (t Tree) Write(w io.Writer, depth []bool) {
	if depth == nil { // root
		fmt.Fprintf(w, "%s\n", t.Task.Info())
	}
	for n, task := range t.Tasks {
		last := n+1 == len(t.Tasks)
		t.fill(w, depth)
		if last {
			fmt.Fprintf(w, lastItem)
		} else {
			fmt.Fprintf(w, middleItem)
		}
		fmt.Fprintf(w, "%s\n", task.Info())
		if len(task.Pomodoros) > 0 && t.ShowPomodoros {
			t.fill(w, depth)
			if last {
				fmt.Fprintf(w, emptySpace)
			} else {
				fmt.Fprintf(w, continueItem)
			}
			fmt.Fprintf(w, "%s*", format.TruncDuration(task.Duration))
			for _, p := range task.Pomodoros {
				fmt.Fprintf(w, "%s", p.Info(task.Duration))
			}

			fmt.Fprintf(w, "\n")
		}
		next := Tree{
			Task:          *task,
			ShowPomodoros: t.ShowPomodoros,
		}
		next.Write(w, t.next(len(task.Tasks) > 0 && !last, depth))
	}
}

func (t Tree) String() string {
	buf := bytes.NewBuffer(nil)
	t.Write(buf, nil)
	return buf.String()
}
