package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

const (
	botName = "sheilaon7_bot"
)

func main() {
	twitterClient, err := NewTwitterClient()
	panicOnError(err)

	files, err := ioutil.ReadFile("lyrics.txt")
	panicOnError(err)

	lyrics := strings.Split(string(files), "---")
	for i := range lyrics {
		lyrics[i] = strings.Trim(lyrics[i], "\n")
	}

	var lyric string

	for {
		seed := rand.NewSource(time.Now().UnixNano())
		random := rand.New(seed)
		index := random.Intn(len(lyrics))

		lyric = lyrics[index]
		fmt.Println("New tweet:", lyric)

		user, _, err := twitterClient.Users.Show(&twitter.UserShowParams{
			ScreenName: botName,
		})
		panicOnError(err)

		if user.Status != nil {
			lastTweet, _, err := twitterClient.Statuses.Show(user.Status.ID, nil)
			panicOnError(err)

			lastTweetSentences := strings.Split(lastTweet.Text, "\n")
			lyricSentences := strings.Split(lyric, "\n")

			if lastTweetSentences[0] == lyricSentences[0] {
				fmt.Println("New tweet shouldn't be same as previous tweet")
				fmt.Println()
				continue
			}
		}

		break
	}

	lyric += " #sheilaon7 #sheilagank"

	_, _, err = twitterClient.Statuses.Update(lyric, nil)
	panicOnError(err)
}

// NewTwitterClient returns a new TwitterClient.
func NewTwitterClient() (*twitter.Client, error) {
	twitterConsumerAPIKey := os.Getenv("TWITTER_CONSUMER_API_KEY")
	if twitterConsumerAPIKey == "" {
		return nil, errors.New("$TWITTER_CONSUMER_API_KEY must be set")
	}

	twitterConsumerAPISecret := os.Getenv("TWITTER_CONSUMER_API_SECRET")
	if twitterConsumerAPISecret == "" {
		return nil, errors.New("$TWITTER_CONSUMER_API_SECRET must be set")
	}

	twitterAccessToken := os.Getenv("TWITTER_ACCESS_TOKEN")
	if twitterAccessToken == "" {
		return nil, errors.New("$TWITTER_ACCESS_TOKEN must be set")
	}

	twitterAccessSecret := os.Getenv("TWITTER_ACCESS_SECRET")
	if twitterAccessSecret == "" {
		return nil, errors.New("$TWITTER_ACCESS_SECRET must be set")
	}

	oauthConfig := oauth1.NewConfig(twitterConsumerAPIKey, twitterConsumerAPISecret)
	oauthToken := oauth1.NewToken(twitterAccessToken, twitterAccessSecret)
	oauthClient := oauthConfig.Client(oauth1.NoContext, oauthToken)

	return twitter.NewClient(oauthClient), nil
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}
