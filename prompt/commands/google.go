package commands

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

var ErrNothingToGoogle = errors.New("please provide something to google")

type GoogleSearchCommand struct{}

func (c *GoogleSearchCommand) Name() string {
	return "google"
}

func (c *GoogleSearchCommand) Usage() string {
	return "search for something on google. Example: google how to cook"
}

func (c *GoogleSearchCommand) Execute(args []string) (string, error) {
	if len(args) == 0 {
		return "", NewAgentError(ErrNothingToGoogle)
	}

	log.Default().Println("Google search:", args[0])

	url := fmt.Sprintf("https://www.google.com/search?q=%s", url.QueryEscape(args[0]))

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("search query failed with status code %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	linksAndTitle := make([]map[string]string, 0)

	doc.Find("div.g").Each(func(i int, s *goquery.Selection) {
		link := make(map[string]string)

		title := s.Find("h3").Text()
		link["title"] = title

		href, exists := s.Find("a").Attr("href")
		if exists {
			link["url"] = href
		}

		linksAndTitle = append(linksAndTitle, link)
	})

	if len(linksAndTitle) == 0 {
		return "", fmt.Errorf("no results found")
	}

	ret := ""
	for i, link := range linksAndTitle {
		ret += fmt.Sprintf("Link #%d: %s\nTitle: %s\n\n", i+1, link["url"], link["title"])
	}

	return ret, nil
}

func init() {
	DefaultCommands = append(DefaultCommands, &GoogleSearchCommand{})
}
