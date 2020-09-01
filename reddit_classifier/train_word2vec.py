import csv, ast
from gensim.models import Word2Vec
from utils import *


with open("data/reddit.csv", "r") as file:
    data = list(csv.reader(file, delimiter=','))


data = data[1:]
sentences, labels = transform_data(data)
print(len(sentences))
for i in range(100):
    print(i, sentences[i], labels[i])

sentences = [elem.lower().split(' ') for elem in sentences]
sentences = [list(filter(None, sequence)) for sequence in sentences]
print(sentences[0:10])

model = Word2Vec(sentences, size=100, window=5, min_count=1, workers=4,iter=5)
model.save("word2vec_reddit.model")
