package main

func main() {

	a := App{}

	a.Initialize("root", "password", "todo")

	a.Run(":8080")
}
