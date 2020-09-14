
--
create keyspace reddit_storage with replication = { 'class' : 'NetworkTopologyStrategy', 'datacenter1' : 3 };
ALTER KEYSPACE reddit_storage WITH replication = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };
ALTER KEYSPACE reddit_storage WITH replication = {'class': 'NetworkTopologyStrategy', 'DC1': '1'};

create table reddit_storage.comments(
    body text,
    id text,
    name text,
    subreddit text,
    subreddit_id text,
    link_id text,
    link_title text,
    likes int,
    permalink text,
    ups int,
    downs int,
    author text,
    author_fullname text,
    replies text,
    parent_id text,
    created double,
    created_utc double,
    num_comments int,
    PRIMARY KEY(name)
    );


create table reddit_storage.submissions(
    title text,
    id text,
    name text,
    subreddit text,
    subreddit_id text,
    likes int,
    permalink text,
    url text,
    ups int,
    downs int,
    author text,
    author_fullname text,
    created double,
    created_utc double,
    num_comments int,
    PRIMARY KEY(name)
    );


create table reddit_storage.classifications(
   name text,
   proba_hateful double,
   proba_not_hateful double,
   class int,
   PRIMARY KEY(name)
);


--
-- drop table reddit_storage.submissions;
-- drop table reddit_storage.comments;
curl -XPUT -H 'Content-Type: application/json' 'http://localhost:9200/comments' -d '{
    "settings": { "keyspace":"reddit_storage" },
    "mappings": {
        "comments" : {
            "discover":".*"
        }
    }
}'


curl -XGET "http://localhost:9200/comments/_search?pretty=true" \
-H "Content-Type: application/json" \
-d '
{
    "query": {
        "bool": {
            "must": [
                {
                    "term": {"parent_id" : "t1_g57m23q"}
                }
            ]
        }
    }
}
'


curl -XPUT "http://localhost:9200/submissions" \
-H "Content-Type: application/json" \
-d '{
    "settings": {
        "keyspace": "reddit_storage"
    },
    "mapping": {
        "submissions": {
            "_doc": ".*"
        }
    }
}'


curl -XGET http://localhost:9200/_cat/indices

curl -XPUT "http://localhost:9200/classifications" \
-H "Content-Type: application/json" \
-d '{
  "settings":{
    "keyspace": "reddit_storage"
  },
  "mappings": {
    "classifications" : {
      "discover" : ".*"
    }
  }
}'




