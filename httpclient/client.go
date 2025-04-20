package httpclient

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

const NotKnownGood = "answer does not match known-good answer; there is a regression in your code"

var (
	matchRe = regexp.MustCompile(`Your puzzle answer was <code>"?([a-zA-Z0-9, ]+)"?</code>\.`)
)

type Client struct {
	cache        cache
	sessionToken string
}

func NewClient() *Client {
	sessionToken := os.Getenv("SESSION_TOKEN")

	if sessionToken == "" {
		fmt.Println("no SESSION_TOKEN env var")
		os.Exit(1)
	}

	return &Client{
		cache:        cache{},
		sessionToken: sessionToken,
	}
}

func (c *Client) GetPreviouslySubmittedSolution(year, day, part int) (string, error) {
	rubricUrl := fmt.Sprintf("https://adventofcode.com/%d/day/%d", year, day)
	body, err := c.get(rubricUrl)
	if err != nil {
		return "", err
	}

	lines := strings.Split(body, "\n")
	var solution string

	numFound := 0

	for _, line := range lines {
		matches := matchRe.FindStringSubmatch(line)
		if len(matches) > 1 {
			solution = matches[1]
			numFound++
			if part == numFound {
				return solution, nil
			}
		}
	}

	return "", nil
}

func (c *Client) GetInput(year, day int) (string, error) {
	cachedInput, found, err := c.cache.GetInput(year, day)
	if err != nil {
		return "", err
	}

	if found {
		return cachedInput, nil
	}

	inputUrl := fmt.Sprintf("https://adventofcode.com/%d/day/%d/input", year, day)
	body, err := c.get(inputUrl)
	if err != nil {
		return "", err
	}

	err = c.cache.PutInput(year, day, body)
	if err != nil {
		return "", fmt.Errorf("caching input: %w", err)
	}

	return body, nil
}

func (c *Client) HasSolution(year, day, part int) bool {
	_, found := c.cache.GetSolution(year, day, part)

	return found
}

func (c *Client) SubmitAnswer(year, day, part int, answer string) (string, error) {
	if data, found := c.cache.GetSolution(year, day, part); found {
		if data == answer {
			return "That's the right answer!", nil
		}

		return NotKnownGood, nil
	}

	found, err := c.cache.GetAttempt(year, day, part, answer)
	if err != nil {
		return "", fmt.Errorf("checking for previous attempt: %w", err)
	}

	if found {
		return "", fmt.Errorf("answer '%s' is known to be incorrect. Not submitting again", answer)
	}

	previousSolution, err := c.GetPreviouslySubmittedSolution(year, day, part)
	if err != nil {
		return "", fmt.Errorf("getting previous solution: %w", err)
	}

	if previousSolution != "" {
		err := c.cache.PutSolution(year, day, part, previousSolution)
		if err != nil {
			return "", fmt.Errorf("caching solution: %w", err)
		}

		if answer == previousSolution {
			return "That's the right answer!", nil
		}
	}

	answerUrl := fmt.Sprintf("https://adventofcode.com/%d/day/%d/answer", year, day)

	formData := url.Values{}
	formData.Set("level", fmt.Sprintf("%d", part))
	formData.Set("answer", answer)

	fmt.Printf("submitting for day %d, year %d, part %d\n", day, year, part)

	response := c.post(answerUrl, bytes.NewBufferString(formData.Encode()))
	fmt.Printf(response)

	return response, nil
}

func (c *Client) KnownBadAnswer(year, day, part int, answer string) error {
	return c.cache.PutAttempt(year, day, part, answer)
}

func (c *Client) get(url string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.Fatalf("making request: %s", err)
	}

	sessionCookie := http.Cookie{
		Name:  "session",
		Value: c.sessionToken,
	}
	req.AddCookie(&sessionCookie)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("making request: %s", err)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("reading response body: %s", err)
	}

	// specific error message from AOC site. Cache should avoid this...
	if strings.HasPrefix(string(body), "Please don't repeatedly") {
		return "", fmt.Errorf("repeated request")
	}

	return string(body), nil
}

func (c *Client) post(url string, body *bytes.Buffer) string {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "POST", url, body)
	if err != nil {
		log.Fatalf("making request: %s", err)
	}

	sessionCookie := http.Cookie{
		Name:  "session",
		Value: c.sessionToken,
	}
	req.AddCookie(&sessionCookie)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("making request: %s", err)
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("reading response body: %s", err)
	}

	if strings.HasPrefix(string(resBody), "Please don't repeatedly") {
		log.Fatalf("Repeated request error")
	}

	return string(resBody)
}
