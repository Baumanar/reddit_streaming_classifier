package api_models

type SubmissionListing struct {
	Kind string                `json:"kind"`
	Data SubmissionListingData `json:"data"`
}

type SubmissionListingData struct {
	Modhash  interface{}        `json:"modhash"`
	Dist     int64              `json:"dist"`
	Children []PostListingChild `json:"children"`
	After    string             `json:"after"`
	Before   interface{}        `json:"before"`
}

type PostListingChild struct {
	Kind string     `json:"kind"`
	Data Submission `json:"data"`
}

func (p *SubmissionListing) UnwrapData() ([]Submission, error) {
	var comments []Submission
	for _, children := range p.Data.Children {
		comments = append(comments, children.Data)
	}
	return comments, nil
}

type Submission struct {
	Title             string  `json:"title" cql:"title"`
	ID                string  `json:"id" cql:"id"`
	Name              string  `json:"name" cql:"name"`
	Subreddit         string  `json:"subreddit" cql:"subreddit"`
	SubredditID       string  `json:"subreddit_id" cql:"subreddit_id"`
	Likes             int64   `json:"likes" cql:"likes"`
	Permalink         string  `json:"permalink" cql:"permalink"`
	URL               string  `json:"url" cql:"url"`
	Ups               int64   `json:"ups" cql:"ups"`
	Downs             int64   `json:"downs" cql:"downs"`
	Author            string  `json:"author" cql:"author"`
	AuthorFullname    string  `json:"author_fullname" cql:"author_fullname"`
	Created           float64 `json:"created" cql:"created"`
	CreatedUTC        float64 `json:"created_utc" cql:"created_utc"`
	NumComments       int64   `json:"num_comments" cql:"num_comments"`
	IsClassified      bool    `json:"is_classified" cql: "is_classified"`
	IsHatespeech      bool    `json:"is_hatespech" cql: "is_hatespech"`
	ProbaHateful      float64 `json:"proba_hateful" cql:"proba_hateful"`
	ProbaNotHateful   float64 `json:"proba_not_hateful" cql:"proba_not_hateful"`
	RemovedByCategory string  `json:"removed_by_category"`
	//Selftext                   string        `json:"selftext"`
	//Saved                      bool          `json:"saved"`
	//Gilded                     int64         `json:"gilded"`
	//Clicked                    bool          `json:"clicked"`
	//SubredditNamePrefixed      string        `json:"subreddit_name_prefixed"`
	//Hidden                     bool          `json:"hidden"`
	//Pwls                       int64         `json:"pwls"`
	//ThumbnailHeight            int64         `json:"thumbnail_height"`
	//HideScore                  bool          `json:"hide_score"`
	//Quarantine                 bool          `json:"quarantine"`
	//LinkFlairTextColor         string        `json:"link_flair_text_color"`
	//UpvoteRatio                int64         `json:"upvote_ratio"`
	//AuthorFlairBackgroundColor string        `json:"author_flair_background_color"`
	//SubredditType              string        `json:"subreddit_type"`
	//TotalAwardsReceived        int64         `json:"total_awards_received"`
	//MediaEmbed                 Gildings      `json:"media_embed"`
	//ThumbnailWidth             int64         `json:"thumbnail_width"`
	//IsOriginalContent          bool          `json:"is_original_content"`
	//IsRedditMediaDomain        bool          `json:"is_reddit_media_domain"`
	//IsMeta                     bool          `json:"is_meta"`
	//SecureMediaEmbed           Gildings      `json:"secure_media_embed"`
	//CanModPost                 bool          `json:"can_mod_post"`
	//Score                      int64         `json:"score"`
	//AuthorPremium              bool          `json:"author_premium"`
	//Thumbnail                  string        `json:"thumbnail"`
	//Edited                     bool          `json:"edited"`
	//AuthorFlairCSSClass        string        `json:"author_flair_css_class"`
	//Gildings                   Gildings      `json:"gildings"`
	//PostHint                   string        `json:"post_hint"`
	//IsSelf                     bool          `json:"is_self"`
	//Wls                        int64         `json:"wls"`
	//AuthorFlairType            string        `json:"author_flair_type"`
	//Domain                     string        `json:"domain"`
	//AllowLiveComments          bool          `json:"allow_live_comments"`
	//URLOverriddenByDest        string        `json:"url_overridden_by_dest"`
	//ViewCount                  int64   `json:"view_count"`
	//Archived                   bool          `json:"archived"`
	//NoFollow                   bool          `json:"no_follow"`
	//IsCrosspostable            bool          `json:"is_crosspostable"`
	//Pinned                     bool          `json:"pinned"`
	//Over18                     bool          `json:"over_18"`
	//Preview                    Preview       `json:"preview"`
	//MediaOnly                  bool          `json:"media_only"`
	//CanGild                    bool          `json:"can_gild"`
	//Spoiler                    bool          `json:"spoiler"`
	//Locked                     bool          `json:"locked"`
	//AuthorFlairText            string        `json:"author_flair_text"`
	//Visited                    bool          `json:"visited"`
	//LinkFlairBackgroundColor   string        `json:"link_flair_background_color"`
	//IsRobotIndexable           bool          `json:"is_robot_indexable"`
	//SendReplies                bool          `json:"send_replies"`
	//WhitelistStatus            string        `json:"whitelist_status"`
	//ContestMode                bool          `json:"contest_mode"`
	//AuthorPatreonFlair         bool          `json:"author_patreon_flair"`
	//AuthorFlairTextColor       string        `json:"author_flair_text_color"`
	//ParentWhitelistStatus      string        `json:"parent_whitelist_status"`
	//Stickied                   bool          `json:"stickied"`
	//SubredditSubscribers       int64         `json:"subreddit_subscribers"`
	//NumCrossposts              int64         `json:"num_crossposts"`
	//IsVideo                    bool          `json:"is_video"`
}
