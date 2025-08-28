package todo

import (
	"os"
	"testing"
)

func TestListAdd(t *testing.T) {
	tests := []struct {
		name     string
		l        List
		taskName string
		want     string
	}{
		{
			name:     "add a new task",
			l:        List{},
			taskName: "New Task",
			want:     "New Task",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.l.Add(test.taskName)

			if test.want != test.l[0].Task {
				t.Errorf("expected %v, got %v", test.want, test.l[0].Task)
			}
		})
	}
}

func TestListComplete(t *testing.T) {
	tests := []struct {
		name     string
		l        List
		taskName string
		wantName string
	}{
		{
			name:     "add a new task",
			l:        List{},
			taskName: "New Task",
			wantName: "New Task",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.l.Add(test.taskName)
			if test.wantName != test.l[0].Task {
				t.Errorf("expected %v, got %v", test.wantName, test.l[0].Task)
			}
			if test.l[0].Done == true {
				t.Errorf("new task should not be completed")
			}

			_ = test.l.Complete(1)
			if test.l[0].Done != true {
				t.Errorf("new task should be completed")
			}
		})
	}
}

func TestListDelete(t *testing.T) {
	tests := []struct {
		name  string
		l     List
		tasks []string
	}{
		{
			name: "delete the second task",
			l:    List{},
			tasks: []string{
				"New Task 1",
				"New Task 2",
				"New Task 3",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for _, v := range test.tasks {
				test.l.Add(v)
			}

			if test.tasks[0] != test.l[0].Task {
				t.Errorf("expected l[0].Task = %v, got %v", test.tasks[0], test.l[0].Task)
			}

			_ = test.l.Delete(2)
			if len(test.l) != 2 {
				t.Errorf("expected len(l) = %v, got %v", 2, len(test.l))
			}
			if test.tasks[2] != test.l[1].Task {
				t.Errorf("expected l[1].Task = %v after deleting, got %v", test.tasks[2], test.l[1].Task)
			}
		})
	}
}

func TestListSaveGet(t *testing.T) {
	tests := []struct {
		name     string
		l1       List
		l2       List
		taskName string
	}{
		{
			name:     "save task to a file and get it",
			l1:       List{},
			l2:       List{},
			taskName: "New Task",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.l1.Add(test.taskName)
			if test.taskName != test.l1[0].Task {
				t.Errorf("expected l1[0].Task = %v, got %v", test.taskName, test.l1[0].Task)
			}

			testFile, err := os.CreateTemp("", "")
			if err != nil {
				t.Fatalf("Error creating temp file: %s", err)
			}
			defer func() { _ = os.Remove(testFile.Name()) }()

			if err = test.l1.Save(testFile.Name()); err != nil {
				t.Fatalf("Error saving list to file: %s", err)
			}

			if err = test.l2.Get(testFile.Name()); err != nil {
				t.Fatalf("Error geting list from file: %s", err)
			}

			if test.l1[0].Task != test.l2[0].Task {
				t.Errorf("Task %v should match task %v", test.l1[0].Task, test.l2[0].Task)
			}
		})
	}
}
