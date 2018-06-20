package twitter

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/dutchcoders/marija/server/datasources"
	"github.com/gernest/mention"
)

var (
	_ = datasources.Register("twitter-tweets", NewTweets)
)

func NewTweets(options ...func(datasources.Index) error) (datasources.Index, error) {
	s := TwitterTweets{}

	for _, optionFn := range options {
		optionFn(&s)
	}

	config := oauth1.NewConfig(s.ConsumerKey, s.ConsumerSecret)

	httpClient := config.Client(
		oauth1.NoContext,
		oauth1.NewToken(s.Token, s.TokenSecret),
	)

	s.client = twitter.NewClient(httpClient)

	return &s, nil
}

func (m *TwitterTweets) Type() string {
	return "twitter"
}

type TwitterTweets struct {
	Config

	client *twitter.Client
	u      *url.URL
}

func (m *TwitterTweets) UnmarshalTOML(p interface{}) error {
	data, _ := p.(map[string]interface{})

	if v, ok := data["consumer_key"]; !ok {
	} else if v, ok := v.(string); !ok {
	} else {
		m.ConsumerKey = v
	}

	if v, ok := data["consumer_secret"]; !ok {
	} else if v, ok := v.(string); !ok {
	} else {
		m.ConsumerSecret = v
	}

	if v, ok := data["token"]; !ok {
	} else if v, ok := v.(string); !ok {
	} else {
		m.Token = v
	}

	if v, ok := data["token_secret"]; !ok {
	} else if v, ok := v.(string); !ok {
	} else {
		m.TokenSecret = v
	}

	return nil
}

func (i *TwitterTweets) Search(ctx context.Context, so datasources.SearchOptions) datasources.SearchResponse {
	itemCh := make(chan datasources.Item)
	errorCh := make(chan error)

	go func() {
		defer close(itemCh)
		defer close(errorCh)

		stp := twitter.SearchTweetParams{
			Count: 200,
		}

		if strings.HasPrefix(so.Query, "geo:") {
			parts := strings.Split(so.Query, " ")

			stp.Geocode = parts[0][4:]
			stp.Query = strings.Join(parts[1:], " ")
		} else {
			stp.Query = so.Query
		}

		search, _, err := i.client.Search.Tweets(&stp)

		if err != nil {
			errorCh <- err
			return
		}

		for _, tweet := range search.Statuses {
			fields := map[string]interface{}{
				"text": tweet.Text,
				"in_reply_to_screen_name":   tweet.InReplyToScreenName,
				"in_reply_to_status_id_str": tweet.InReplyToStatusIDStr,
				"in_reply_to_user_id_str":   tweet.InReplyToUserIDStr,
				"lang":               tweet.Lang,
				"source":             tweet.Source,
				"user.name":          tweet.User.Name,
				"user.id_str":        tweet.User.IDStr,
				"user.lang":          tweet.User.Lang,
				"user.location":      tweet.User.Location,
				"user.screen_name":   tweet.User.ScreenName,
				"user.profile_image": tweet.User.ProfileImageURLHttps,
				"tags":               mention.GetTags('#', strings.NewReader(tweet.Text), ':'),
				"mentions":           mention.GetTags('@', strings.NewReader(tweet.Text), ':'),
			}

			if tweet.Coordinates == nil {
			} else if tweet.Coordinates.Type != "Point" {
			} else {
				coordinates := fmt.Sprintf("%f,%f", tweet.Coordinates.Coordinates[1], tweet.Coordinates.Coordinates[0])
				fields["coordinates"] = coordinates
			}

			item := datasources.Item{
				ID:     tweet.IDStr,
				Fields: fields,
			}

			select {
			case itemCh <- item:
			case <-ctx.Done():
				return
			}
		}

	}()

	return datasources.NewSearchResponse(
		itemCh,
		errorCh,
	)
}

func (i *TwitterTweets) GetFields(ctx context.Context) (fields []datasources.Field, err error) {
	fields = append(fields, datasources.Field{
		Path: "text",
		Type: "string",
	})

	fields = append(fields, datasources.Field{
		Path: "in_reply_to_screen_name",
		Type: "string",
	})

	fields = append(fields, datasources.Field{
		Path: "in_reply_to_user_id",
		Type: "string",
	})

	fields = append(fields, datasources.Field{
		Path: "in_reply_to_status_id",
		Type: "string",
	})

	fields = append(fields, datasources.Field{
		Path: "created_at",
		Type: "string",
	})

	fields = append(fields, datasources.Field{
		Path: "user.screen_name",
		Type: "string",
	})

	fields = append(fields, datasources.Field{
		Path: "user.name",
		Type: "string",
	})

	fields = append(fields, datasources.Field{
		Path: "user.id_str",
		Type: "string",
	})

	fields = append(fields, datasources.Field{
		Path: "mentions",
		Type: "string",
	})

	fields = append(fields, datasources.Field{
		Path: "tags",
		Type: "string",
	})

	fields = append(fields, datasources.Field{
		Path: "coordinates",
		Type: "location",
	})

	fields = append(fields, datasources.Field{
		Path: "user.profile_image",
		Type: "image",
	})

	return
}
