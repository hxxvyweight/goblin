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
	fmt.Println("|> create\n|> q [quit]\n|> help\n")
}
