package api_models

type CommentListingList []CommentListing

type CommentListing struct {
	Kind string             `json:"kind"`
	Data CommentListingData `json:"data"`
}

func (p *CommentListing) UnwrapData() ([]Comment, error) {
	var comments []Comment
	for _, children := range p.Data.Children {
		comments = append(comments, children.Data)
	}
	return comments, nil
}

type CommentListingData struct {
	Modhash  string         `json:"modhash"`
	Dist     float64        `json:"dist"`
	Children []CommentChild `json:"children"`
	After    string         `json:"after"`
	Before   string         `json:"before"`
}

type CommentChild struct {
	Kind string  `json:"kind"`
	Data Comment `json:"data"`
}

type Comment struct {
	Body            string  `json:"body" cql:"body"`
	ID              string  `json:"id" cql:"id"`
	Name            string  `json:"name" cql:"name"`
	Subreddit       string  `json:"subreddit" cql:"subreddit"`
	SubredditID     string  `json:"subreddit_id" cql:"subreddit_id"`
	LinkID          string  `json:"link_id" cql:"link_id"`
	LinkTitle       string  `json:"link_title" cql:"link_title"`
	Likes           int64   `json:"likes" cql:"likes"`
	Permalink       string  `json:"permalink" cql:"permalink"`
	Ups             int64   `json:"ups" cql:"ups"`
	Downs           int64   `json:"downs" cql:"downs"`
	Author          string  `json:"author" cql:"author"`
	AuthorFullname  string  `json:"author_fullname" cql:"author_fullname"`
	Replies         string  `json:"replies" cql:"replies"`
	ParentID        string  `json:"parent_id" cql:"parent_id"`
	Created         float64 `json:"created" cql:"created"`
	CreatedUTC      float64 `json:"created_utc" cql:"created_utc"`
	NumComments     int64   `json:"num_comments" cql:"num_comments"`
	IsClassified    bool    `json:"is_classified" cql: "is_classified"`
	IsHatespeech    bool    `json:"is_hatespech" cql: "is_hatespech"`
	ProbaHateful    float64 `json:"proba_hateful" cql:"proba_hateful"`
	ProbaNotHateful float64 `json:"proba_not_hateful" cql:"proba_not_hateful"`
	//LinkPermalink                string        `json:"link_permalink"`
	//TotalAwardsReceived          int64         `json:"total_awards_received"`
	//ApprovedAtUTC                interface{}   `json:"approved_at_utc"`
	//Edited                       bool          `json:"edited"`
	//AuthorFlairType              string        `json:"author_flair_type"`
	//Saved                        bool          `json:"saved"`
	//Gilded                       int64         `json:"gilded"`
	//Archived                     bool          `json:"archived"`
	//NoFollow                     bool          `json:"no_follow"`
	//CanModPost                   float64          `json:"can_mod_post"`
	//SendReplies                  bool          `json:"send_replies"`
	//Score                        int64         `json:"score"`
	//Over18                       bool          `json:"over_18"`
	//AuthorPatreonFlair           bool          `json:"author_patreon_flair"`
	//IsSubmitter                  bool          `json:"is_submitter"`
	//BodyHTML                     string        `json:"body_html"`
	//Gildings                     Gildings      `json:"gildings"`
	//Stickied                     bool          `json:"stickied"`
	//AuthorPremium                bool          `json:"author_premium"`
	//CanGild                      bool          `json:"can_gild"`
	//SubredditNamePrefixed        string        `json:"subreddit_name_prefixed"`
	//ScoreHidden                  bool          `json:"score_hidden"`
	//LinkAuthor                   string        `json:"link_author"`
	//LinkURL                      string        `json:"link_url"`
	//Collapsed                    bool          `json:"collapsed"`
	//Controversiality             int64         `json:"controversiality"`
	//Locked                       bool          `json:"locked"`
	//Quarantine                   bool          `json:"quarantine"`
	//SubredditType                string        `json:"subreddit_type"`
}
