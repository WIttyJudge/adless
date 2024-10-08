package main

import "fmt"

var (
	version   = "unknown"
	gitCommit = "unknown"
	buildDate = "unknown"
)

func main() {
	fmt.Println(version)
	fmt.Println(gitCommit)
	fmt.Println(buildDate)
}
