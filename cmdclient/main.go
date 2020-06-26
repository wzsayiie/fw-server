//cmdclient

package main

import (
	"fmt"
)

func input() string {
	var cmd string = ""

	fmt.Printf("input >> ")
	fmt.Scanf("%s", &cmd)

	return cmd
}

func run(cmd string) {
	fmt.Printf("receive: %s\n", cmd)
}

func main() {
	for {
		var cmd string = input()
		run(cmd)
	}
}
