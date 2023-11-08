package main

type User struct {
	ID       int
	Name     string
	Priority int
}

func NewUser(name string, pri int) *User {
	return &User{
		Name:     name,
		Priority: pri,
	}
}

type Task struct {
	Name        string
	Description string
	AssignedBy  *User
	AssignedTo  *User
	SubTasks    []*Task
}

func NewTask(name, desc string) *Task {
	return &Task{
		Name:        name,
		Description: desc,
	}
}
