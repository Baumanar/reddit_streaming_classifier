
--
create keyspace reddit_storage with replication = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };

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

