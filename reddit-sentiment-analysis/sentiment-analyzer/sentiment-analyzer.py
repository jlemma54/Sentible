from operator import index
from turtle import pd
from pandas.core.base import DataError
from pandas.core.frame import DataFrame
import requests
import pandas
import json
import os
from nltk.sentiment.vader import SentimentIntensityAnalyzer
from nltk.tokenize import RegexpTokenizer
import en_core_web_sm
from nltk.stem import WordNetLemmatizer



def load_data(path) -> pandas.DataFrame: 
   df = pandas.read_csv(path)
   dropped_columns = ['all_awardings', 'associated_award', 'author_flair_background_color', 'author_flair_css_class', 'author_flair_template_id','awarders', 'author_flair_richtext', 'author_flair_text', 
    'author_flair_text_color', 'collapsed_because_crowd_control', 'gildings', 'no_follow', 'score', 'send_replies', 'stickied', 'treatment_tags', 'edited', 'top_awarded_type', 'author_cakeday', 'distinguished']
   for col in dropped_columns: 
       try: 
           df.drop(col, axis=1, inplace=True)
       except: 
           pass
   return df.fillna('')


def calculate_sentiment(df: pandas.DataFrame, vader, nlp, stopwords, tokenizer, lemmatizer, url): 
    try:
        print('here')
        comments = df.iloc[df['analyzed'].tolist().index(False):]
        for comment in comments.iterrows():
            r = df[df['id1'] == comment[1][13]].index.values[0].item()
            if requests.get(url + comment[1][13]).status_code != 200:
                tokenized_string = tokenizer.tokenize(comment[1][7])
                lower_tokenized = [word.lower() for word in tokenized_string] 
                sw_removed = [word for word in lower_tokenized if not word in stopwords]
                lemmatized_tokens = ([lemmatizer.lemmatize(w) for w in sw_removed])
                sentiment_dict = vader.polarity_scores(' '.join(lemmatized_tokens))
                
                print(sentiment_dict)

                df.loc[r, 'positive'] = sentiment_dict['pos']
                df.loc[r, 'negative'] = sentiment_dict['neg']
                df.loc[r, 'neutral'] = sentiment_dict['neu']
                df.loc[r, 'compound'] = sentiment_dict['compound']
                df.loc[r, 'analyzed'] = True
                df.to_csv(os.environ.get('DATAFRAME'), index=False)
            else:
                df.loc[r, 'analyzed'] = True
                df.to_csv(os.environ.get('DATAFRAME'), index=False) 
                print("comment already in database")
            
    except: 
        print('failed')
        pass




def update_db(df: pandas.DataFrame, url): 
    payload_df = df.loc[df['analyzed']== True]

    epoch = [ 
    dict([
    (colname, row[i]) 
    for i,colname in enumerate(payload_df.columns)
    ])
    for row in df.values
    ] 
    
    for i in range(len(epoch)): 
       z = requests.post(url=url, json=epoch[i]) 
       print(z.text)
    cleanup_csv(os.environ.get('DATAFRAME'))
    
def cleanup_csv(csv_path): 
    df = load_data(csv_path)
    df = df[df['analyzed']==False]
    df.to_csv(csv_path, index=False)
    
    


df = load_data(os.environ.get('DATAFRAME'))
nlp = en_core_web_sm.load()
vader = SentimentIntensityAnalyzer()
stopwords = nlp.Defaults.stop_words
tokenizer = RegexpTokenizer('\w+|\$[\d\.]+|http\S+')
lemmatizer = WordNetLemmatizer()
calculate_sentiment(df=df, vader=vader, nlp=nlp, stopwords=stopwords, tokenizer=tokenizer, lemmatizer=lemmatizer, url=os.environ.get('URL')+'/')
df = load_data(os.environ.get('DATAFRAME'))
update_db(df=df, url=os.environ.get('URL'))







