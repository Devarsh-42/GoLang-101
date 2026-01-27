package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type CmdFlags struct {
	Add string
	Del int
	Edit string
	Toggle int
	List bool
}


func NewCmdFlags() *CmdFlags {
	cf := CmdFlags{}
	flag.StringVar(&cf.Add, "add", "", "Add a new Todo Item") // flag.StringVar(pointer to variable, flag name, default value, description)
	flag.IntVar(&cf.Del, "del", -1, "Delete a Todo Item by index")
	flag.StringVar(&cf.Edit, "edit", "", "Edit a Todo Item by index")
	flag.IntVar(&cf.Toggle, "toggle", -1, "Toggle a Todo Item by index")
	flag.BoolVar(&cf.List, "list" , false, "List all Todo Items")
	flag.Parse()
	return &cf

}

func (cf *CmdFlags) Execute(todos *Todos) {
	switch {
	case cf.Add != "":
		todos.add(cf.Add)
	case cf.Del != -1:
			todos.delete(cf.Del)
	case cf.Edit != "":
		parts := strings.SplitN(cf.Edit, ":", 2)
		if len(parts) != 2 {
			fmt.Println("Invalid format for edit. Use -edit=index:new title")
			os.Exit(1)
		}
		index, err := strconv.Atoi(parts[0])
		if err != nil {
			fmt.Println("Invalid index for edit.")
			os.Exit(1)
		}
		newTitle := parts[1]
		todos.edit(index, newTitle)
	case cf.Toggle != -1:
		todos.toggle(cf.Toggle)
	case cf.List:
		todos.print()

	default:
		fmt.Println("No valid command provided. Use -h for help.")
	}
	
} 