package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/devusSs/minls/internal/cli"
)

var (
	buildVersion   string
	buildDate      string
	buildGitCommit string
)

func init() {
	if buildVersion == "" {
		buildVersion = "development"
	}

	if buildDate == "" {
		buildDate = "unknown"
	}

	if buildGitCommit == "" {
		buildGitCommit = "unknown"
	}
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "version" {
		printVersion()
		os.Exit(0)
	}

	if len(os.Args) > 1 && os.Args[1] == "help" {
		printHelp()
		os.Exit(0)
	}

	if len(os.Args) > 1 {
		code := handleCommandLine()
		os.Exit(code)
	}

	// this may be replaced by a TUI mode later
	printNoCommandHelp("")
	os.Exit(1)
}

func printVersion() {
	printHeader()
	fmt.Printf("Build version:\t\t%s\n", buildVersion)
	fmt.Printf("Build date:\t\t%s\n", buildDate)
	fmt.Printf("Build Git commit:\t%s\n", buildGitCommit)
	fmt.Println()
	fmt.Printf("Build Go version:\t%s\n", runtime.Version())
	fmt.Printf("Build Go OS:\t\t%s\n", runtime.GOOS)
	fmt.Printf("Build Go arch:\t\t%s\n", runtime.GOARCH)
}

func handleCommandLine() int {
	// we can be sure len(os.Args) > 1
	// because we check that before
	// this function is invoked
	command := os.Args[1]
	switch command {
	case "list":
		fmt.Println("list command, not implemented")
		return 0
	case "upload":
		err := cli.Upload()
		if err != nil {
			fmt.Println("UPLOAD FAILED:", err)
			return 1
		}
		return 0
	case "download":
		fmt.Println("download command, not implemented")
		return 0
	case "delete":
		fmt.Println("delete command, not implemented")
		return 0
	case "clear":
		fmt.Println("clear command, not implemented")
		return 0
	default:
		printNoCommandHelp(command)
		return 1
	}
}

func printNoCommandHelp(command string) {
	if command == "" {
		fmt.Println("ERROR: no command specified")
	} else {
		fmt.Printf("ERROR: unknown command specified: '%s'\n", command)
	}
	fmt.Println()
	printHelp()
}

func printHelp() {
	printHeader()
	fmt.Println("Usage:")
	fmt.Println("	minls version				Prints version / build information")
	fmt.Println("	minls help				Prints this help message")
	fmt.Println(
		"	minls list				Prints all uploaded files if possible and available",
	)
	fmt.Println(
		"	minls upload <filepath> <policy>	Uploads the specified file (policies: private / public)",
	)
	fmt.Println(
		"	minls download <id> <filepath>		Downloads the specified file to the specified filepath",
	)
	fmt.Println("	minls delete <id>			Deletes the specified file")
	fmt.Println(
		"	minls clear <option>			Clears program data (options: all / data / logs / downloads)",
	)
}

func printHeader() {
	fmt.Println("minls - Go tool to combine MinIO and YOURLS")
	fmt.Println()
	fmt.Println("github.com/devusSs/minls")
	fmt.Println()
}
