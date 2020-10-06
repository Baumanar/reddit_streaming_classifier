## Elassandra initialization and request examples

Run the following command to get a cmd line to the kubernetes elassandra pod:

```
kubectl exec -it -n default $(kubectl get pods -n default -l app=elassandra,release=my-release-elassandra -o jsonpath='{.items[0].metadata.name}') -- bin/bash
```

Once connected, you can create indexes

#### Create index on comments
```
curl -XPUT -H 'Content-Type: application/json' 'http://localhost:9200/comments' -d '{
    "settings": { "keyspace":"reddit_storage" },
    "mappings": {
        "comments" : {
            "discover":".*"
        }
    }
}'
```


#### Create index on submissions

```
curl -XPUT -H 'Content-Type: application/json' 'http://localhost:9200/submissions' -d '{
    "settings": { "keyspace":"reddit_storage" },
    "mappings": {
        "submissions" : {
            "discover":".*"
        }
    }
}'
```

## Request examples


#### Get hateful comments

```
curl -XGET "http://localhost:9200/comments/_search?pretty=true" \
-H "Content-Type: application/json" \
-d '
{
    "query": {
        "bool": {
            "must": [
                {
                    "term": {"is_hatespeech" : "true"}
                }
            ]
        }
    }
}
'
