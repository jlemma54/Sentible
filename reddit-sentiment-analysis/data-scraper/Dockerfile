FROM python:3

RUN mkdir /app

WORKDIR /app

ENV DATAFRAME /app/stock.csv

ENV BATCHES /app/batches.csv

ADD data.py /

ADD data-scraper.py /

RUN touch ${DATAFRAME}

RUN touch ${BATCHES}

RUN pip instal numpy

RUN pip install anytree

RUN pip install pmaw

RUN pip install pandas

CMD ["python", "./data-scraper.py"]
