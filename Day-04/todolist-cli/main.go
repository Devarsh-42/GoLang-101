package main


func main() {
	todos := Todos{}
	storage := NewStorage[Todos]("todos.json")
	storage.Load(&todos)

    todos.toggle(0)

    CmdFlags := NewCmdFlags()
    CmdFlags.Execute(&todos)

    // todos.print()

    storage.Save(todos)

	// todos.add("Buy Milk")
	// todos.add("Buy Bread")
	// fmt.Printf("%+v\n\n", todos)
	// todos.delete(0)
	// fmt.Printf("%+v", todos)
	
}
