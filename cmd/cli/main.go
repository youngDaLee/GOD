package cli

import (
	"fmt"
	"os"
)

func RunCLI() {
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "init":
		handleInit()
	case "start":
		handleStart()
	case "analyze":
		handleAnalyze()
	case "config":
		handleConfig()
	case "version":
		handleVersion()
	case "help":
		printHelp()
	case "god":
		handleGod()
	default:
		fmt.Fprintf(os.Stderr, "알 수 없는 명령어: %s\n", os.Args[1])
		printHelp()
		os.Exit(1)
	}
}

func handleInit() {
	fmt.Println("init")
}

func handleStart() {
	fmt.Println("start")
}

func handleAnalyze() {
	fmt.Println("analyze")
}

func handleConfig() {
	fmt.Println("config")
}

func handleVersion() {
	fmt.Println("version")
}

func handleGod() {
	fmt.Println("god! god!! god!!!")
}

func printHelp() {
	fmt.Println("how to use god cli")
}
