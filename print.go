package main

import "fmt"

const TermText = "===================================================================="

func PrintGoblin() {

	fmt.Println(`
          /\_/\                    _______
         ( o.o )       .-.        /       \
         / > ^ \      /   \      |  $ $ $  |
        / / \ \ \    | L O O T!  |_________|
       (_/   \_)      \___/      [___________]
      /           \              | $$$$$$$$$ |
     /    GOBLIN   \             | $$$$$$$$$ |
    /_______________\            |___________|`)

	fmt.Println(">>> GOBLIN CLIENT INITIALIZED!")
	fmt.Println(">>> Ready to venture into the archives...\n")
}
func ByeGoblin() {
	fmt.Println(`
        /|\ Bagging the loot and heading out...

         ______
        /      \
       |  LOOT  |
        \______/
           | \
           |  \      /\_/\
           |   \    ( o.o )---> *strolls into the sunset*
            \        > ^ <
             \      / / \ \
              \    (_/   \_)
               \             \
                \   GOBLIN    \
                 \______________\
                 `)
}

func PrintHelp() {
	fmt.Println("Commands:\n create | delete | push | pull | start | stop | q [quit]")

}

func PrintInstancePrompt() {
	fmt.Print("\\|/ goblin |:: Specify a instance e.g: [1] || [5]\n")
	fmt.Print("\\|/ goblin |:: ")
}
func PrintFilePrompt() {
	fmt.Print("\\|/ goblin |:: Specify a file e.g: example.zip \n")
	fmt.Print("\\|/ goblin |:: ")
}
