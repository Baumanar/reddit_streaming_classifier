package data

import (
	"errors"
	"github.com/Baumanar/reddit_streaming_classifier/reddit_kafka/api_models"
	"github.com/gocql/gocql"
	"log"
	"strings"
)

func CreateComment(comment *api_models.Comment, session *gocql.Session) error {
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

func CreateSubmission(submission *api_models.Submission, session *gocql.Session) error {
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


func CreateClassification(classification *Classification, session *gocql.Session) error {
	q := `
		INSERT INTO classifications(
		    name,
			proba_hateful,
			proba_not_hateful,
		    class
		)
		VALUES (?, ?, ?, ?)
    	`
	err := session.Query(q,
		classification.Name,
		classification.ProbaHateful,
		classification.ProbaNotHateful,
		classification.Class,
	).Exec()
	if err != nil {
		log.Printf("ERROR: fail create comment, %s", err.Error())
	}

	return err
}

func GetComment(name string, session *gocql.Session) (*api_models.Comment, error) {

	m := map[string]interface{}{}
	q := `
		SELECT * FROM comments
			WHERE name = ?
		LIMIT 1
    	`
	itr := session.Query(q, name).Consistency(gocql.One).Iter()
	for itr.MapScan(m) {
		comment := &api_models.Comment{}
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

func GetSubmission(name string, session *gocql.Session) (*api_models.Submission, error) {

	m := map[string]interface{}{}
	q := `
		SELECT * FROM submissions
			WHERE name = ?
		LIMIT 1
    	`
	itr := session.Query(q, name).Consistency(gocql.One).Iter()
	for itr.MapScan(m) {
		submission := &api_models.Submission{}
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

func UpdateComment(name string, params map[string]interface{}, session *gocql.Session) error {

	q := `UPDATE comments `
	values := make([]interface{}, 0)
	for k, v := range params {
		q += "SET " + k + "= ?, "
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

func UpdateSubmission(name string, params map[string]interface{}, session *gocql.Session) error {

	q := `UPDATE submissions `
	values := make([]interface{}, 0)
	for k, v := range params {
		q += "SET " + k + "= ?, "
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
