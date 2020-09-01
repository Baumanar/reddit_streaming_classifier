# The Reddit Streaming classifier

This project aim is to classify directly streamed reddit comments from some subreddits on some criteria.
In a first version this project detects hate speech comments from new comments or title of submission.

The project is written in Golang for the streaming components and the storage utilities and Python for 
the classification. The services communicate with Kafka and the data is stored in a cassandra database.


## Dependencies

Go: 
- [kafka](https://github.com/confluentinc/confluent-kafka-go)
- [cassandra driver](https://github.com/gocql/gocql)
- [reddit streaming](https://github.com/Baumanar/reddit_api_streaming)

Python:

- Keras
- Kafka

And others: [see python requirements](https://github.com/Baumanar/reddit_proj/blob/master/reddit_classifier/requirements.txt)

## Build & Run instructions:

Setup the environnemnt in `reddit_classifier`:

- Create a virtual env: `python3 -m venv venv`
- Activate the virtual env: `. venv/bin/activate`
- Install the required packages : `pip3 install -r requirements.txt`

- Train the word2vec model: `python train_word2vec.py`
- Train the hate speech classifier: `python train_lstm.py`


To run everything together:

- Build the images: `docker-compose -f docker-compose.yml build`
- Run the docker-compose: `docker-compose -f docker-compose.yml up`