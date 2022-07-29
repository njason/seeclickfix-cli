package main

import (
	"encoding/json"
	"encoding/csv"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

type IssuesResp struct {
	Issues	[]Issue	`json:"issues"`
	Metadata Metadata `json:"metadata"`
}

type Metadata struct {
	Pagination struct {
		Entries         int         `json:"entries"`
		Page            int         `json:"page"`
		PerPage         int         `json:"per_page"`
		Pages           int         `json:"pages"`
		NextPage        *int				`json:"next_page"`
		NextPageURL     *string			`json:"next_page_url"`
		PreviousPage		*int				`json:"previous_page"`
		PreviousPageURL *string			`json:"previous_page_url"`
	} `json:"pagination"`
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
    category := reportCmd.String("category", "", "ex. trees")

    if len(os.Args) < 2 {
        fmt.Println("expected 'report' subcommand")
        os.Exit(1)
    }

    switch os.Args[1] {
		case "report":
			reportCmd.Parse(os.Args[2:])
			report(*placeUrl, *category);
		default:
			fmt.Println("expected 'report' subcommand")
			os.Exit(1)
    }
}

func report(placeUrl string, category string) {
	var issues = filterIssues(category, issuesRequest(placeUrl, 1))
	csvData := csvFormat(issues)

	csvWriter := csv.NewWriter(os.Stdout)
	csvWriter.WriteAll(csvData) // calls Flush internally

	if err := csvWriter.Error(); err != nil {
		fmt.Println("error writing csv:", err)
	}
}

func mapCategoryToOrganization(category string) string {
	switch category {
	case "trees":
		return "Trees"
	default:
		fmt.Println("expected 'trees' category")
		os.Exit(1)
		return ""  // why, go, why
	}
}

func filterIssues(category string, issues []Issue) []Issue {
	var organization = mapCategoryToOrganization(category)
	var filteredIssues []Issue

	for _, issue := range issues {
		if issue.RequestType.Organization == organization {
			filteredIssues = append(filteredIssues, issue)
		}
	}

	return filteredIssues
}

func issuesRequest(placeUrl string, page int) []Issue {
	req, _ := http.NewRequest("GET", "https://seeclickfix.com/api/v2/issues", nil)
	query := req.URL.Query()
	query.Add("place_url", placeUrl)
	query.Add("page", strconv.FormatInt(int64(page), 10))
	req.URL.RawQuery = query.Encode()

	rawResp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer rawResp.Body.Close()

	var resp IssuesResp
	err = json.NewDecoder(rawResp.Body).Decode(&resp)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	// paginate recursively
	if resp.Metadata.Pagination.NextPage != nil {
		// rate limit is 20 request per minute.
		time.Sleep(3 * time.Second)

		return append(resp.Issues, issuesRequest(placeUrl, page+1)...)
	}

	return resp.Issues
}

func csvFormat(issues []Issue) [][]string {
	csvData := [][]string{
		{"summary", "status", "created_at", "link", "lat", "lng"},
	}

	for _, issue := range issues {
		lat := fmt.Sprintf("%f", issue.Point.Coordinates[1])
		lng := fmt.Sprintf("%f", issue.Point.Coordinates[0])
		record := []string{issue.Summary, issue.CreatedAt, issue.HTMLURL, lat, lng}
		csvData = append(csvData, record)
	}

	return csvData
}
