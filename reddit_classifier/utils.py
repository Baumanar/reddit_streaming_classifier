import ast
import re
import numpy as np


def process_lines(data):
    for idx, row in enumerate(data):
        data[idx][1] = row[1].split('\n')
        data[idx][2] = row[2].split('\n')
    return data


def transform_labels(data):
    for i in range(len(data)):
        # for j in range(len(data[i])):
        if data[i][2][0] == "n/a":
            data[i][2] = []
        else:
            data[i][2] = ast.literal_eval(data[i][2][0])
        if data[i][3] == "n/a":
            data[i][3] = []
        else:
            data[i][3] = ast.literal_eval(data[i][3])
    return data


def sanitize_sentences(sentences):
    for idx, sentence in enumerate(sentences):
        sentence = sentence.lower()
        sentence = re.sub('^[0-9]+\.', '', sentence)
        sentence = re.sub('[^A-Za-z0-9 ]+', '', sentence)

        sentences[idx] = sentence
    return sentences


def transform_data(data):
    data = process_lines(data)
    data = transform_labels(data)
    all_sentences = []
    all_labels = []
    for i, discussion in enumerate(data):
        went_ok = True
        labels = [0 for i in range(len(discussion[1]))]
        for idx in discussion[2]:
            try:
                labels[idx - 1] = 1
            except:
                went_ok = False
        if went_ok:
            sentences = sanitize_sentences(sentences)

            # sentences = [e.lower() for e in string if e.isalnum() for string in sentences]
            all_labels += labels
            all_sentences += sentences
    return all_sentences, all_labels


def word2token(word, w2v_model):
    try:
        return w2v_model.wv.vocab[word].index
    # If word is not in index return 0. I realize this means that this
    # is the same as the word of index 0 (i.e. most frequent word), but 0s
    # will be padded later anyway by the embedding layer (which also
    # seems dirty but I couldn't find a better solution right now)
    except KeyError:
        return 0


def token2word(token, w2v_model):
    return w2v_model.wv.index2word[token]


MAX_SEQUENCE_LENGTH = 200


def tokenize(sentences, w2v_model):
    tokenized_sentences = []
    for sentence in sentences:
        # Make all characters lower-case
        words = np.array([word2token(w, w2v_model) for w in sentence.split(' ')[:MAX_SEQUENCE_LENGTH] if w != ''])
        tokenized_sentences.append(words)
    return tokenized_sentences
