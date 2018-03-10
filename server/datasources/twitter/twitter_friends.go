package twitter

import (
	"context"
	"net/url"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/dutchcoders/marija/server/datasources"
)

var (
	_ = datasources.Register("twitter-friends", NewFriends)
)

func NewFriends(options ...func(datasources.Index) error) (datasources.Index, error) {
	s := TwitterFriends{}

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

func (m *TwitterFriends) Type() string {
	return "twitter-friends"
}

type TwitterFriends struct {
	Config

	client *twitter.Client
	u      *url.URL
}

func (m *TwitterFriends) UnmarshalTOML(p interface{}) error {
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

func (i *TwitterFriends) Search(ctx context.Context, so datasources.SearchOptions) datasources.SearchResponse {
	itemCh := make(chan datasources.Item)
	errorCh := make(chan error)

	go func() {
		defer close(itemCh)
		defer close(errorCh)

		cursor := int64(0)

		for {
			friends, _, err := i.client.Friends.List(&twitter.FriendListParams{
				Cursor:     cursor,
				ScreenName: so.Query,
			})

			if err != nil {
				errorCh <- err
				return
			}

			cursor = friends.NextCursor

			for _, user := range friends.Users {
				fields := map[string]interface{}{
					"user.screen_name":   so.Query,
					"friend.screen_name": user.ScreenName,
					"friend.email":       user.Email,
					"friend.name":        user.Name,
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

func (i *TwitterFriends) GetFields(ctx context.Context) (fields []datasources.Field, err error) {
	fields = append(fields, datasources.Field{
		Path: "user.screen_name",
		Type: "string",
	})

	fields = append(fields, datasources.Field{
		Path: "friend.email",
		Type: "string",
	})

	fields = append(fields, datasources.Field{
		Path: "friend.name",
		Type: "string",
	})

	fields = append(fields, datasources.Field{
		Path: "friend.screen_name",
		Type: "string",
	})

	return
}
