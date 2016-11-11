package server

type Index interface {
	Search(query string) ([]Item, error)
	// Indexes() ([]Index, error)
	// Fields(string) ([]Fields, error)
}
