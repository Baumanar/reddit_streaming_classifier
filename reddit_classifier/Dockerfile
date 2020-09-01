
FROM tensorflow/tensorflow:latest

WORKDIR .
COPY . .
RUN pip install -r requirements.txt


ENTRYPOINT [ "python", "reddit_classifier.py"]