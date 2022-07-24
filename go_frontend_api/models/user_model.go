package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Comment struct {
	Id                    primitive.ObjectID `json:"id,omitempty"`
	Author                string             `json:"author,omitempty" validate:"required"`
	Author_Flair_Type     string             `json:"author_flair_type,omitempty" validate:"required"`
	Author_Fullname       string             `json:"author_fullname,omitempty" validate:"required"`
	Author_Patreon_Flair  bool               `json:"author_patreon_flair,omitempty"`
	Author_Premium        bool               `json:"author_premium,omitempty"`
	Body                  string             `json:"body,omitempty" validate:"required"`
	Created_Utc           int64              `json:"created_utc,omitempty" validate:"required"`
	Id1                   string             `json:"id1,omitempty" validate:"required"`
	Is_Submitter          bool               `json:"is_submitter,omitempty"`
	Link_Id               string             `json:"link_id,omitempty" validate:"required"`
	Locked                bool               `json:"locked,omitempty"`
	Parent_Id             string             `json:"parent_id,omitempty" validate:"required"`
	Permalink             string             `json:"permalink,omitempty" validate:"required"`
	Retrieved_On          int64              `json:"retrieved_on,omitempty" validate:"required"`
	Subreddit             string             `json:"subreddit,omitempty" validate:"required"`
	Subreddit_Id          string             `json:"subreddit_id,omitempty" validate:"required"`
	Total_Awards_Received int64              `json:"total_awards_received,omitempty"`
	Stock                 string             `json:"stock,omitempty"`
	Positive_Sentiment    float32            `json:"positive_sentiment,omitempty"`
	Negative_Sentiment    float32            `json:"negative_sentiment,omitempty"`
	Neutral_Sentiment     float32            `json:"neutral_sentiment,omitempty"`
	Compound_Sentiment    float32            `json:"compound_sentiment,omitempty"`
	Analyzed              bool               `json:"analyzed,omitempty"`
}
