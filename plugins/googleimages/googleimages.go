package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/disiqueira/ultraslackbot/pkg/bot"
	"github.com/disiqueira/ultraslackbot/pkg/plugin"
	"github.com/disiqueira/ultraslackbot/pkg/slack"

	"net/url"
)

const (
	pattern          = "(?i)\\b(googleimages|gimages|image|gis)\\b"
	googleKeyEnvName = "GOOGLE_KEY"
	googleCXEnvName  = "GOOGLE_CX"
	searchURL 		 = "https://www.googleapis.com/customsearch/v1"
)

type (
	googleImages struct {
		plugin.BasicCommand
		cx string
		googleKey string
	}

	googleImagesResponse struct {
		Items []struct {
			Title            string `json:"title"`
			Link             string `json:"link"`
		} `json:"items"`
	}
)

func (gi *googleImages) Start() error {
	gi.googleKey = os.Getenv(googleKeyEnvName)
	if gi.googleKey == "" {
		return fmt.Errorf("enviroiment variable %s not found", googleKeyEnvName)
	}

	gi.cx = os.Getenv(googleCXEnvName)
	if gi.cx == "" {
		return fmt.Errorf("enviroiment variable %s not found", googleCXEnvName)
	}

	return nil
}

func (gi *googleImages) Name() string {
	return "googleImages"
}

func (gi *googleImages) Execute(event slack.Event, botUser bot.UserInfo) ([]slack.Message, error) {
	return gi.HandleEvent(event, botUser, gi.matcher, gi.command)
}

func (gi *googleImages) matcher() *regexp.Regexp {
	return regexp.MustCompile(pattern)
}

func (gi *googleImages) command(text string) (string, error) {
	args := strings.Split(strings.TrimSpace(text), " ")
	if len(args) < 2 {
		return "", nil
	}

	text = strings.Join(args[1:], " ")

	var gisURL *url.URL
	gisURL, err := url.Parse(searchURL)
	if err != nil {
		return "", err
	}

	parameters := url.Values{}
	parameters.Add("key",gi.googleKey)
	parameters.Add("cx",gi.cx)
	parameters.Add("searchType", "image")
	parameters.Add("q", text)
	gisURL.RawQuery = parameters.Encode()

	data := &googleImagesResponse{}
	if err := plugin.GetJSON(gisURL.String(), data); err != nil {
		return "", err
	}

	title, link := "", ""
	for _, item := range data.Items {
		if len(item.Link) > 0 {
			title = item.Title
			link = item.Link
			break
		}
	}

	return fmt.Sprintf("%s - %s", title, link), nil
}

var CustomPlugin googleImages
