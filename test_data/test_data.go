package test_data

import (
	metadata "github.com/supersherm5/movie-app/metadata/pkg/model"
	rating "github.com/supersherm5/movie-app/rating/pkg/model"
)

// FakeMetadata for memory repository
var FakeMetaData = map[string]*metadata.Metadata{
	"1": {
		ID:          "1",
		Title:       "The Shawshank Redemption",
		Description: "Two imprisoned",
		Director:    "Frank Darabont",
	},
	"2": {
		ID:          "2",
		Title:       "The Godfather",
		Description: "The aging patriarch of an organized crime dynasty transfers control of his clandestine empire to his reluctant son.",
		Director:    "Francis Ford Coppola",
	},
	"3": {
		ID:          "3",
		Title:       "The Dark Knight",
		Description: "When the menace known as the Joker wreaks havoc and chaos on the people of Gotham, Batman must accept one of the greatest psychological and physical tests of his ability to fight injustice.",
		Director:    "Christopher Nolan",
	},
	"4": {
		ID:          "4",
		Title:       "12 Angry",
		Description: "A jury holdout attempts to prevent a miscarriage of justice by forcing his colleagues to reconsider the evidence.",
		Director:    "Sidney Lumet",
	},
}

var FakeRatingData = map[rating.RecordType]map[rating.RecordID][]rating.Rating{
	rating.RecordTypeMovie: {
		"1": {
			{"1", rating.RecordTypeMovie, "1", 5},
			{"1", rating.RecordTypeMovie, "2", 4},
			{"1", rating.RecordTypeMovie, "3", 3},
		},
		"2": {
			{"1", rating.RecordTypeMovie, "1", 5},
			{"1", rating.RecordTypeMovie, "2", 4},
			{"1", rating.RecordTypeMovie, "3", 3},
		},
		"3": {
			{"1", rating.RecordTypeMovie, "1", 5},
			{"1", rating.RecordTypeMovie, "2", 4},
			{"1", rating.RecordTypeMovie, "3", 3},
		},
		"4": {
			{"1", rating.RecordTypeMovie, "1", 5},
			{"1", rating.RecordTypeMovie, "2", 4},
			{"1", rating.RecordTypeMovie, "3", 3},
		},
	},
}
