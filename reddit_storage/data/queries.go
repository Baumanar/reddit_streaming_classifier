package data

import (
	"errors"
	"github.com/Baumanar/reddit_streaming_classifier/reddit_kafka/models"
	"github.com/gocql/gocql"
	"log"
	"strings"
)

// CreateComment creates a Comment in the database, without classif info
func CreateComment(comment *models.Comment, session *gocql.Session) error {
	q := `
		INSERT INTO comments (
		    body,
		    id,
		    name,
    		subreddit,
    		subreddit_id,
		    link_id,
			link_title,
		    likes,
			permalink,
			ups,
			downs,
			author,
			author_fullname,
			replies,
			parent_id,
			created,
			created_utc, 
			num_comments
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    	`
	err := session.Query(q,
		comment.Body,
		comment.ID,
		comment.Name,
		comment.Subreddit,
		comment.SubredditID,
		comment.LinkID,
		comment.LinkTitle,
		comment.Likes,
		comment.Permalink,
		comment.Ups,
		comment.Downs,
		comment.Author,
		comment.AuthorFullname,
		comment.Replies,
		comment.ParentID,
		comment.Created,
		comment.CreatedUTC,
		comment.NumComments).Exec()
	if err != nil {
		log.Printf("ERROR: fail create comment, %s", err.Error())
	}

	return err
}

// CreateSubmission creates a Submission in the database, without classif info
func CreateSubmission(submission *models.Submission, session *gocql.Session) error {
	q := `
		INSERT INTO submissions (
			title,
		    id,
		    name,
			subreddit,
			subreddit_id,
		    likes,
		    permalink,
			url,
			ups,
			downs,
			author,
			author_fullname,
			created,
			created_utc,
			num_comments
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    	`
	err := session.Query(q,
		submission.Title,
		submission.ID,
		submission.Name,
		submission.Subreddit,
		submission.SubredditID,
		submission.Likes,
		submission.Permalink,
		submission.URL,
		submission.Ups,
		submission.Downs,
		submission.Author,
		submission.AuthorFullname,
		submission.Created,
		submission.CreatedUTC,
		submission.NumComments).Exec()
	if err != nil {
		log.Printf("ERROR: fail create comment, %s", err.Error())
	}

	return err
}

// GetComment returns a Comment by its full name
func GetComment(name string, session *gocql.Session) (*models.Comment, error) {

	m := map[string]interface{}{}
	q := `
		SELECT * FROM comments
			WHERE name = ?
		LIMIT 1
    	`
	itr := session.Query(q, name).Consistency(gocql.One).Iter()
	for itr.MapScan(m) {
		comment := &models.Comment{}
		comment.Body = m["body"].(string)
		comment.ID = m["id"].(string)
		comment.Name = m["name"].(string)
		comment.LinkID = m["link_id"].(string)
		comment.LinkTitle = m["link_title"].(string)
		comment.Likes = m["likes"].(int64)
		comment.Permalink = m["permalink"].(string)
		comment.Ups = m["ups"].(int64)
		comment.Downs = m["downs"].(int64)
		comment.Author = m["author"].(string)
		comment.AuthorFullname = m["author_fullname"].(string)
		comment.Replies = m["replies"].(string)
		comment.ParentID = m["parent_id"].(string)
		comment.Created = m["created"].(float64)
		comment.CreatedUTC = m["created"].(float64)
		comment.NumComments = m["created"].(int64)
		log.Printf("INFO: found comment, %v", comment)

		return comment, nil
	}

	return nil, errors.New("document not found")
}

// GetSubmission returns a Submission by its full name
func GetSubmission(name string, session *gocql.Session) (*models.Submission, error) {

	m := map[string]interface{}{}
	q := `
		SELECT * FROM submissions
			WHERE name = ?
		LIMIT 1
    	`
	itr := session.Query(q, name).Consistency(gocql.One).Iter()
	for itr.MapScan(m) {
		submission := &models.Submission{}
		submission.Title = m["title"].(string)
		submission.ID = m["id"].(string)
		submission.Name = m["name"].(string)
		submission.Subreddit = m["name"].(string)
		submission.SubredditID = m["name"].(string)
		submission.Likes = m["likes"].(int64)
		submission.Permalink = m["permalink"].(string)
		submission.URL = m["url"].(string)
		submission.Ups = m["ups"].(int64)
		submission.Downs = m["downs"].(int64)
		submission.Author = m["author"].(string)
		submission.AuthorFullname = m["author_fullname"].(string)
		submission.Created = m["created"].(float64)
		submission.CreatedUTC = m["created_utc"].(float64)
		submission.NumComments = m["num_comments"].(int64)
		log.Printf("INFO: found comment, %v", submission)

		return submission, nil
	}

	return nil, errors.New("submission not found")
}

// UpdateComment updates a Comment
func UpdateComment(name string, params map[string]interface{}, session *gocql.Session) error {

	q := `UPDATE comments SET `
	values := make([]interface{}, 0)
	for k, v := range params {
		q += k + "= ?, "
		values = append(values, v)
	}

	q = strings.TrimSuffix(q, ", ")
	q += " WHERE name = ?"
	values = append(values, name)

	err := session.Query(q, values...).Exec()
	if err != nil {
		log.Printf("ERROR: fail update document, %s", err.Error())
		return err
	}

	return nil
}

// UpdateSubmission updates a Submission
func UpdateSubmission(name string, params map[string]interface{}, session *gocql.Session) error {

	q := `UPDATE submissions SET `
	values := make([]interface{}, 0)
	for k, v := range params {
		q += k + "= ?, "
		values = append(values, v)
	}

	q = strings.TrimSuffix(q, ", ")
	q += " WHERE name = ?"
	values = append(values, name)

	err := session.Query(q, values...).Exec()
	if err != nil {
		log.Printf("ERROR: fail update document, %s", err.Error())
		return err
	}

	return nil
}

// UpdateCommentClassification updates a Comment classification info
func UpdateCommentClassification(name string, isClassified bool, isHatespeech bool, probaHateful float64, probaNotHateful float64, session *gocql.Session) error {
	params := map[string]interface{}{
		"is_classified":     isClassified,
		"is_hatespeech":     isHatespeech,
		"proba_hateful":     probaHateful,
		"proba_not_hateful": probaNotHateful,
	}
	err := UpdateComment(name, params, session)
	if err != nil {
		log.Printf("ERROR: fail UpdateCommentClassification, %s", err.Error())
		return err
	}
	return nil
}

// UpdateSubmissionClassification updates a Submission classification info
func UpdateSubmissionClassification(name string, isClassified bool, isHatespeech bool, probaHateful float64, probaNotHateful float64, session *gocql.Session) error {
	params := map[string]interface{}{
		"is_classified":     isClassified,
		"is_hatespeech":     isHatespeech,
		"proba_hateful":     probaHateful,
		"proba_not_hateful": probaNotHateful,
	}
	err := UpdateSubmission(name, params, session)
	if err != nil {
		log.Printf("ERROR: fail UpdateSubmissionClassification, %s", err.Error())
		return err
	}
	return nil
}

// UpdateClassification updates a record classification info according to its type tp
func UpdateClassification(tp string, name string, isClassified bool, isHatespeech bool, probaHateful float64, probaNotHateful float64, session *gocql.Session) error {
	if tp == "comment" {
		err := UpdateCommentClassification(name, isClassified, isHatespeech, probaHateful, probaNotHateful, session)
		return err
	}
	if tp == "submission" {
		err := UpdateSubmissionClassification(name, isClassified, isHatespeech, probaHateful, probaNotHateful, session)
		return err
	}
	return errors.New("unknown type")
}
