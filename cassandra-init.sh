CQL="create keyspace IF NOT EXISTS reddit_storage with replication = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };

create table IF NOT EXISTS reddit_storage.comments(
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
    is_classified boolean,
    is_hatespeech boolean,
    proba_hateful double,
    proba_not_hateful double,
    PRIMARY KEY(name)
    );


create table IF NOT EXISTS reddit_storage.submissions(
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
    is_classified boolean,
    is_hatespeech boolean,
    proba_hateful double,
    proba_not_hateful double,
    PRIMARY KEY(name)
    );
"

until echo $CQL | cqlsh; do
  echo "cqlsh: Cassandra is unavailable to initialize - will retry later"
  sleep 2
done &

exec /docker-entrypoint.sh "$@"