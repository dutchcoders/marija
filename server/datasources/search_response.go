package datasources

type SearchResponse interface {
	Item() chan Item
	Error() chan error
}

func NewSearchResponse(itemChan chan Item, errorChan chan error) *searchResponse {
	return &searchResponse{
		itemChan,
		errorChan,
	}
}

type searchResponse struct {
	itemChan  chan Item
	errorChan chan error
}

func (sr *searchResponse) Item() chan Item {
	return sr.itemChan
}

func (sr *searchResponse) Error() chan error {
	return sr.errorChan
}
