package models

// Gilding ...
type Gilding struct {
	Gid map[string]int `json:"gid"`
}

// ResizedIcon ...
type ResizedIcon struct {
	URL    string `json:"url"`
	Width  int64  `json:"width"`
	Height int64  `json:"height"`
}

// FlairRichtext ...
type FlairRichtext struct {
	E string `json:"e"`
	T string `json:"t"`
}

// Gildings ...
type Gildings struct {
	Gid1 *int64 `json:"gid_1,omitempty"`
	Gid2 *int64 `json:"gid_2,omitempty"`
}

// MediaEmbed ...
type MediaEmbed struct {
}

// Preview ...
type Preview struct {
	Images  []Image `json:"images"`
	Enabled bool    `json:"enabled"`
}

// Image ...
type Image struct {
	Source      ResizedIcon   `json:"source"`
	Resolutions []ResizedIcon `json:"resolutions"`
	Variants    MediaEmbed    `json:"variants"`
	ID          string        `json:"id"`
}

// AwardSubType ...
type AwardSubType string

//const (
//	Global AwardSubType = "GLOBAL"
//)

// AwardType ...
type AwardType string

//const (
//	AwardTypeGlobal AwardType = "global"
//)

// IconFormat ...
type IconFormat string

//const (
//	Apng IconFormat = "APNG"
//	PNG  IconFormat = "PNG"
//)
