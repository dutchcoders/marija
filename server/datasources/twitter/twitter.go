package twitter

import (
	"net/url"
	"strings"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/dutchcoders/marija/server/datasources"
	"github.com/gernest/mention"
	"github.com/kr/pretty"
)

type Twitter struct {
	client *twitter.Client
	u      *url.URL
}

func (m *Twitter) UnmarshalTOML(p interface{}) error {
	data, _ := p.(map[string]interface{})

	consumerKey := ""
	if v, ok := data["consumer_key"]; !ok {
	} else if v, ok := v.(string); !ok {
	} else {
		consumerKey = v
	}

	consumerSecret := ""
	if v, ok := data["consumer_secret"]; !ok {
	} else if v, ok := v.(string); !ok {
	} else {
		consumerSecret = v
	}

	token := ""
	if v, ok := data["token"]; !ok {
	} else if v, ok := v.(string); !ok {
	} else {
		token = v
	}

	tokenSecret := ""
	if v, ok := data["token_secret"]; !ok {
	} else if v, ok := v.(string); !ok {
	} else {
		tokenSecret = v
	}

	config := oauth1.NewConfig(consumerKey, consumerSecret)
	httpClient := config.Client(
		oauth1.NoContext,
		oauth1.NewToken(token, tokenSecret),
	)

	m.client = twitter.NewClient(httpClient)

	return nil
}

func (b *Twitter) Search(so datasources.SearchOptions) ([]datasources.Item, int, error) {
	search, _, err := b.client.Search.Tweets(&twitter.SearchTweetParams{
		Query: so.Query,
		Count: so.Size,
	})

	if err != nil {
		return nil, 0, err
	}

	i := 0

	items := []datasources.Item{}

main:
	for _, tweet := range search.Statuses {

		if i < so.From {
			i++
			continue
		}
		// mentions
		// followers

		fields := map[string]interface{}{
			"text": tweet.Text,
			// "date":             time.Unix(tx.Time, 0),
			"in_reply_to_screen_name":   tweet.InReplyToScreenName,
			"in_reply_to_status_id_str": tweet.InReplyToStatusIDStr,
			"in_reply_to_user_id_str":   tweet.InReplyToUserIDStr,
			"lang":   tweet.Lang,
			"source": tweet.Source,
			"user": map[string]interface{}{
				"name":        tweet.User.Name,
				"id_str":      tweet.User.IDStr,
				"lang":        tweet.User.Lang,
				"location":    tweet.User.Location,
				"screen_name": tweet.User.ScreenName,
			},
			"tags":     mention.GetTags('#', strings.NewReader(tweet.Text), ':'),
			"mentions": mention.GetTags('@', strings.NewReader(tweet.Text), ':'),
		}

		pretty.Print(fields)

		item := datasources.Item{
			ID:     tweet.IDStr,
			Fields: fields,
		}

		items = append(items, item)

		if i >= so.From+so.Size {
			break main
		}

		i++
	}

	return items, search.Metadata.Count, nil
}

func (i *Twitter) Fields() (fields []datasources.Field, err error) {
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
	return
}
