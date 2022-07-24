FROM python:3

ENV URL 127.0.0.1/comment

ADD sentiment-analysis.py /

RUN pip install requests

RUN pip install pandas

RUN pip install nltk

RUN pip install spacy

RUN python -m spacy download en

RUN [ "python3", "-c", "import nltk; nltk.download('vader_lexicon')" ]

RUN [ "python3", "-c", "import nltk; nltk.download('wordnet')" ]

CMD ["python", "./sentiment-analyzer.py"]
