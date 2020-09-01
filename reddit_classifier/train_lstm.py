import collections

from keras.models import Sequential
from keras.layers import *
from keras.preprocessing.sequence import pad_sequences
import numpy as np
import random
from gensim.models import Word2Vec
from utils import *
import csv

w2v_model = Word2Vec.load("word2vec_reddit.model")
# Retrieve the weights from the model. This is used for initializing the weights
# in a Keras Embedding layer later
w2v_weights = w2v_model.wv.vectors
vocab_size, embedding_size = w2v_weights.shape
print(vocab_size, embedding_size)





# Create an iterator that formats data from the dataset proper for
# LSTM training

# Sequences will be padded or truncated to this length
n_categories = 2

with open("data/reddit.csv", "r") as file:
    data = list(csv.reader(file, delimiter=','))

data = data[1:]

sentences, labels = transform_data(data)
print(collections.Counter(labels))

tokenized_sentences = tokenize(sentences, w2v_model)


set_x = pad_sequences(tokenized_sentences, maxlen=MAX_SEQUENCE_LENGTH, padding='pre', value=0)
set_y = np.array(labels)

print(set_x.shape)
print(set_y.shape)
###########################################################################################

VALID_PER = 0.15 # Percentage of the whole set that will be separated for validation

total_samples = set_x.shape[0]
n_val = int(VALID_PER * total_samples)
n_train = total_samples - n_val

random_i = random.sample(range(total_samples), total_samples)
train_x = set_x[random_i[:n_train]]
train_y = set_y[random_i[:n_train]]
val_x = set_x[random_i[n_train:n_train+n_val]]
val_y = set_y[random_i[n_train:n_train+n_val]]

print("Train Shapes - X: {} - Y: {}".format(train_x.shape, train_y.shape))
print("Val Shapes - X: {} - Y: {}".format(val_x.shape, val_y.shape))

###########################################################################################


model = Sequential()

# Keras Embedding layer with Word2Vec weights initialization
model.add(Embedding(input_dim=vocab_size,
                    output_dim=embedding_size,
                    weights=[w2v_weights],
                    input_length=MAX_SEQUENCE_LENGTH,
                    mask_zero=True,
                    trainable=False))

model.add(Bidirectional(LSTM(100)))
model.add(Dense(n_categories, activation='softmax'))

model.compile(optimizer='adam', loss='sparse_categorical_crossentropy', metrics=['accuracy'])

history = model.fit(train_x, train_y, epochs=4, batch_size=32,
                    validation_data=(val_x, val_y), verbose=1)
model.save("lstm_model_keras")



###########################################################################################
