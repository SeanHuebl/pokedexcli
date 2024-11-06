package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	ch := make(chan string)
	go func() {
		defer close(ch)
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			ch <- scanner.Text()
		}
	}()

	for {
		fmt.Print("pokedex > ")
		value := <-ch
		fmt.Println(value)
	}

}
