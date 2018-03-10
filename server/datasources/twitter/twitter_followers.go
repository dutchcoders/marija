package twitter

import (
	"context"
	"net/url"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/dutchcoders/marija/server/datasources"
)

var (
	_ = datasources.Register("twitter-followers", NewFollowers)
)

func NewFollowers(options ...func(datasources.Index) error) (datasources.Index, error) {
	s := TwitterFollowers{}

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

func (m *TwitterFollowers) Type() string {
	return "twitter-followers"
}

type TwitterFollowers struct {
	Config

	client *twitter.Client
	u      *url.URL
}

func (m *TwitterFollowers) UnmarshalTOML(p interface{}) error {
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

func (i *TwitterFollowers) Search(ctx context.Context, so datasources.SearchOptions) datasources.SearchResponse {
	itemCh := make(chan datasources.Item)
	errorCh := make(chan error)

	go func() {
		defer close(itemCh)
		defer close(errorCh)

		cursor := int64(0)

		for {
			followers, _, err := i.client.Followers.List(&twitter.FollowerListParams{
				Cursor:     cursor,
				ScreenName: so.Query,
			})

			if err != nil {
				errorCh <- err
				return
			}

			cursor = followers.NextCursor

			for _, user := range followers.Users {
				fields := map[string]interface{}{
					"user.screen_name":     so.Query,
					"follower.screen_name": user.ScreenName,
					"follower.email":       user.Email,
					"follower.name":        user.Name,
				}

				item := datasources.Item{
					ID:     so.Query + user.ScreenName,
					Fields: fields,
				}

				select {
				case itemCh <- item:
				case <-ctx.Done():
					return
				}
			}

			if cursor == 0 {
				break
			}
		}
	}()

	return datasources.NewSearchResponse(
		itemCh,
		errorCh,
	)
}

func (i *TwitterFollowers) GetFields(ctx context.Context) (fields []datasources.Field, err error) {
	fields = append(fields, datasources.Field{
		Path: "user.screen_name",
		Type: "string",
	})

	fields = append(fields, datasources.Field{
		Path: "follower.email",
		Type: "string",
	})

	fields = append(fields, datasources.Field{
		Path: "follower.name",
		Type: "string",
	})

	fields = append(fields, datasources.Field{
		Path: "follower.screen_name",
		Type: "string",
	})

	return
}
