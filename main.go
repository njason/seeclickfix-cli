package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
)

type IssuesResp struct {
	Issues	[]Issue	`json:"issues"`
}

type Issue struct {
	ID             int         `json:"id"`
	Status         string      `json:"status"`
	Summary        string      `json:"summary"`
	Description    string      `json:"description"`
	Rating         int         `json:"rating"`
	Lat            float64     `json:"lat"`
	Lng            float64     `json:"lng"`
	Address        string      `json:"address"`
	CreatedAt      string      `json:"created_at"`
	AcknowledgedAt interface{} `json:"acknowledged_at"`
	ClosedAt       interface{} `json:"closed_at"`
	ReopenedAt     interface{} `json:"reopened_at"`
	UpdatedAt      string      `json:"updated_at"`
	ShortenedURL   interface{} `json:"shortened_url"`
	URL            string      `json:"url"`
	Point          struct {
		Type        string    `json:"type"`
		Coordinates []float64 `json:"coordinates"`
	} `json:"point"`
	PrivateVisibility bool   `json:"private_visibility"`
	HTMLURL           string `json:"html_url"`
	RequestType       struct {
		ID               int    `json:"id"`
		Title            string `json:"title"`
		Organization     string `json:"organization"`
		URL              string `json:"url"`
		RelatedIssuesURL string `json:"related_issues_url"`
	} `json:"request_type"`
	CommentURL  string `json:"comment_url"`
	FlagURL     string `json:"flag_url"`
	Transitions struct {
		CloseURL string `json:"close_url"`
	} `json:"transitions"`
	Reporter struct {
		ID     int    `json:"id"`
		Name   string `json:"name"`
		Role   string `json:"role"`
		Avatar struct {
			Full          string `json:"full"`
			Square100X100 string `json:"square_100x100"`
		} `json:"avatar"`
		HTMLURL     string `json:"html_url"`
		WittyTitle  string `json:"witty_title"`
		CivicPoints int    `json:"civic_points"`
	} `json:"reporter"`
	Media struct {
		VideoURL               interface{} `json:"video_url"`
		ImageFull              interface{} `json:"image_full"`
		ImageSquare100X100     interface{} `json:"image_square_100x100"`
		RepresentativeImageURL string      `json:"representative_image_url"`
	} `json:"media"`
}

var client = http.DefaultClient

func main() {
    reportCmd := flag.NewFlagSet("report", flag.ExitOnError)
		placeUrl := reportCmd.String("place", "", "Place URL, ex. jersey-city")
    categoryName := reportCmd.String("category", "", "ex. trees")

    if len(os.Args) < 2 {
        fmt.Println("expected 'report' subcommand")
        os.Exit(1)
    }

    switch os.Args[1] {
		case "report":
			reportCmd.Parse(os.Args[2:])
			report(*placeUrl, *categoryName);
		default:
			fmt.Println("expected 'report' subcommand")
			os.Exit(1)
    }
}

func report(placeUrl string, categoryName string) {
	switch categoryName {
	case "trees":
		issues := getTreeRequests(placeUrl)
		fmt.Println(issues)
	default:
		fmt.Println("expected 'trees' category")
		os.Exit(1)
	}
}

func getTreeRequests(placeUrl string) []Issue {
	req, _ := http.NewRequest("GET", "https://seeclickfix.com/api/v2/issues", nil)
	query := req.URL.Query()
	query.Add("place-url", placeUrl)
	req.URL.RawQuery = query.Encode()

	rawResp, _ := client.Do(req)
	defer rawResp.Body.Close()

	var resp IssuesResp
	err := json.NewDecoder(rawResp.Body).Decode(&resp)
	if err != nil {
		panic(err)
	}
	
	return resp.Issues
}
