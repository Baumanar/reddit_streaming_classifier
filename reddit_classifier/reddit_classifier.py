import json
from confluent_kafka import Consumer
from confluent_kafka import Producer
import keras
from keras.preprocessing.sequence import pad_sequences
from utils import *
from gensim.models import Word2Vec
import tensorflow as tf
from sys import argv


tf.get_logger().setLevel('INFO')
tf.autograph.set_verbosity(1)

try:
    bs=argv[1]
    print('\n🥾 bootstrap server: {}'.format(bs))
    bootstrap_server=bs
except:
    # no bs X-D
    bootstrap_server='localhost:9092'
    print('⚠️  No bootstrap server defined, defaulting to {}\n'.format(bootstrap_server))



producer = Producer({'bootstrap.servers': bootstrap_server})

w2v_model = Word2Vec.load("word2vec_reddit.model")
# Retrieve the weights from the model. This is used for initializing the weights
# in a Keras Embedding layer later
w2v_weights = w2v_model.wv.vectors
vocab_size, embedding_size = w2v_weights.shape
model = keras.models.load_model("lstm_model_keras")


MAX_SEQUENCE_LENGTH = 200

consumer = Consumer({
        'bootstrap.servers': bootstrap_server,
        'group.id': 'classifGroup',
        'auto.offset.reset': 'earliest'
    })
consumer.subscribe(["reddit_stream_comments", "reddit_stream_submissions"])
print('Starting')
counter = 0
message_count = 0
sentence_batch = []
batch_ids = []

while True:

    message = consumer.poll(timeout=1.0)
    if message is None:
        continue
    else:
        sentence = json.loads(message.value().decode())
        if message.topic() == "reddit_stream_comments":
            sentence_batch.append(sentence['body'])
            batch_ids.append(sentence['name'])
            message_count += 1
        if message.topic()  == "reddit_stream_submissions":
            sentence_batch.append(sentence['title'])
            batch_ids.append(sentence['name'])
            message_count += 1
        if len(sentence_batch) == 100:

            sentences = sanitize_sentences(sentence_batch)
            sentences = tokenize(sentences, w2v_model)
            set_x = pad_sequences(sentences, maxlen=MAX_SEQUENCE_LENGTH, padding='pre', value=0)
            predictions_proba = model.predict(set_x)
            classes = np.argmax(predictions_proba, axis=1)
            print("    Done pred for batch", counter, len(sentences))
            return_messages = list(zip(batch_ids, predictions_proba, classes))
            for return_message in return_messages:
                json_message = json.dumps({"name":return_message[0],
                                  "proba_hateful":float(return_message[1][1]),
                                  "proba_not_hateful": float(return_message[1][0]),
                                  "class":int(return_message[2])})
                producer.produce("reddit_classification", json_message.encode())
            sentence_batch = []
            batch_ids = []
            counter +=1