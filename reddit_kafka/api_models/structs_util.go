package api_models

type Gilding struct {
	Gid map[string]int `json:"gid"`
}

type ResizedIcon struct {
	URL    string `json:"url"`
	Width  int64  `json:"width"`
	Height int64  `json:"height"`
}

type FlairRichtext struct {
	E string `json:"e"`
	T string `json:"t"`
}

type Gildings struct {
	Gid1 *int64 `json:"gid_1,omitempty"`
	Gid2 *int64 `json:"gid_2,omitempty"`
}

type MediaEmbed struct {
}

type Preview struct {
	Images  []Image `json:"images"`
	Enabled bool    `json:"enabled"`
}

type Image struct {
	Source      ResizedIcon   `json:"source"`
	Resolutions []ResizedIcon `json:"resolutions"`
	Variants    MediaEmbed    `json:"variants"`
	ID          string        `json:"id"`
}

type AwardSubType string

const (
	Global AwardSubType = "GLOBAL"
)

type AwardType string

const (
	AwardTypeGlobal AwardType = "global"
)

type IconFormat string

const (
	Apng IconFormat = "APNG"
	PNG  IconFormat = "PNG"
)
