package main

import (
    "flag"
    "fmt"
    "os"
)

func main() {
    reportCmd := flag.NewFlagSet("report", flag.ExitOnError)
    categoryName := reportCmd.String("category", "", "ex. trees")

    if len(os.Args) < 2 {
        fmt.Println("expected 'report' subcommand")
        os.Exit(1)
    }

    switch os.Args[1] {
		case "report":
			reportCmd.Parse(os.Args[2:])
			report(*categoryName);
		default:
			fmt.Println("expected 'report' subcommand")
			os.Exit(1)
    }
}

func report(categoryName string) {
	switch categoryName {
	case "trees":
		req := getTreeRequests()
		fmt.Println(req)
	default:
		fmt.Println("expected 'trees' category")
		os.Exit(1)
	}
}

func getTreeRequests() int {
	// https://seeclickfix.com/api/v2/issues?place_url=jersey-city
	// "request_type": { "organization": "Trees" }
	return 1
}
