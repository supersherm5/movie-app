package model

// Metadata defines the metadata of a movie.
type Metadata struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Director    string `json:"director"`
}
