from pmaw import PushshiftAPI
import pandas as pd
import os
from anytree import Node, RenderTree, AsciiStyle
from data import *
from numpy import where

pd.set_option("display.max_rows", None, "display.max_columns", None, "display.max_colwidth", None)

api = PushshiftAPI()

def scrape_reddit_posts(before, after, subreddit, limit): 
    
    post_ids = [x['id'] for x in api.search_submissions(subreddit=subreddit, limit=limit, before=before, after=after)] 

    return post_ids


def process_post_id(id): 
    post_ids = [id]
    posts = [post for post in api.search_submissions(ids=post_ids)]
    comment_ids = [c_id for c_id in api.search_submission_comment_ids(ids=post_ids)]
    comments = [comment for comment in api.search_comments(ids=comment_ids)]
    nodes = [Node([x, '']) for x in comment_ids]
    comment_dict = dict(zip(comment_ids, comments))
    head_list = list(set(list(tickers)).intersection(posts[0]['title'].split())) + list(set(list(tickers)).intersection(posts[0]['selftext'].split())) 
    print(posts[0])
    try:
        head = Node([id, '-' + 'link_flair_text:' + posts[0]['link_flair_text'] + '-' + 'tickers: ?' + '?'.join(head_list)])
    except: 
        head = Node([id, '-' + 'link_flair_text:' + 'None' + '-' + 'tickers: ?' + '?'.join(head_list)]) 
    node_names = [y.name[0] for y in nodes]



    for x in nodes: 
        if comment_dict[x.name[0]]['parent_id'][3:]== head.name[0]: 
            x.parent = head
            combined_list = list(set(list(tickers)).intersection(comment_dict[x.name[0]]['body'].split()))
            if len(combined_list) > 0: 
                x.name[1] += head.name[1] + '?' + '?'.join(combined_list)
            else: 
                x.name[1] = head.name[1]
        elif comment_dict[x.name[0]]['parent_id'][3:] in node_names: 
            x.parent = nodes[node_names.index(comment_dict[x.name[0]]['parent_id'][3:])]
            combined_list = list(set(list(tickers)).intersection(comment_dict[x.name[0]]['body'].split()))
            if len(combined_list) > 0: 
                x.name[1] += x.parent.name[1] + '?' + '?'.join(combined_list)
            else: 
                x.name[1] = x.parent.name[1]
    df = pd.DataFrame(comments)
    return head, nodes, df



def update_csv(csv_path, post_id): 
    head, nodes, df = process_post_id(post_id) 
    print(RenderTree(head, style=AsciiStyle()).by_attr(attrname="name"))
    ticker_dict = dict(zip([y.name[0] for y in nodes], [y.name[1] for y in nodes]))

    df.insert(len(df.columns)-1, 'stock', [ticker_dict[id][ticker_dict[id].rfind('?')+1:] for id in [id[1][26] for id in [elem for elem in [sub for sub in [row for row in df.iterrows()]]]]])
    df.insert(len(df.columns)-1, 'link_flair_text', [str(head.name[1])[head.name[1].index(':')+1:head.name[1].index('-', 2)] for i in range(len(df.index))])
    df.insert(len(df.columns)-1, 'positive', ["" for i in range(len(df.index))])
    df.insert(len(df.columns)-1, 'negative', ["" for i in range(len(df.index))])
    df.insert(len(df.columns)-1, 'neutral', ["" for i in range(len(df.index))])
    df.insert(len(df.columns)-1, 'compound', ["" for i in range(len(df.index))])
    df.insert(len(df.columns)-1, 'analyzed', [False for i in range(len(df.index))])


    
    
    dropped_columns = ['all_awardings', 'associated_award', 'author_flair_background_color', 'author_flair_css_class', 'author_flair_template_id','awarders', 'author_flair_richtext', 'author_flair_text', 
    'author_flair_text_color', 'collapsed_because_crowd_control', 'gildings', 'no_follow', 'score', 'send_replies', 'stickied', 'treatment_tags', 'edited', 'top_awarded_type', 'author_cakeday', 'distinguished']
    
    for col in dropped_columns: 
        try:
            df.drop(col, axis=1, inplace=True)
        except:
            pass


    df.rename(columns={'id': 'id1'}, inplace=True)
    df.fillna('')
    print(df['stock'].tolist())
    df.drop(df.index[df['stock'] == ''], inplace=True)
    print(df.columns.values)


    if os.path.isfile(csv_path):
        mode = 'a'
        with open(csv_path, mode=mode) as f:
            df.to_csv(f, header=0, index=False)
    else:
        mode = 'w'
        with open(csv_path, mode=mode) as f:
            df.to_csv(f, header=df.columns.values, index=False)


def retrieve_batches(path): 
    
    """
    should be formatted as such
    before,after,subreddit,limit,completed
    """

    df = pd.read_csv(path)
    rpointer = where(df['completed'] == False)[0][0]
    cpointer = list(df).index('completed')
    temp = df.iloc[rpointer][0].item(), df.iloc[rpointer][1].item(), df.iloc[rpointer][2], df.iloc[rpointer][3].item()
    df.iat[rpointer, cpointer] = True
    df.to_csv(path, mode='w', header=True, index=False, columns=list(df.axes[1]))

    return temp


def main():
    '''
    before, after, subreddit, limit should be initialized by received_batches(), 
    keep it like this for the sake of testing lol
    '''
    while True:
        try:
            before, after, subreddit, limit = retrieve_batches(os.environ.get('BATCHES'))
        except: 
            print("no new batches")
            break
        print(f"before: {before}\nafter: {after}\nsubreddit: {subreddit}\nlimit: {limit}") 
        post_ids = scrape_reddit_posts(before=before, after=after, subreddit=subreddit, limit=limit)
        print(f"{len(post_ids)} post in this epoch...")
        failure_counter = 0
        for post in post_ids: 
            try:
                update_csv(csv_path=os.environ.get('DATAFRAME'), post_id=post)
            except: 
                failure_counter += 1
        print(f"Failed {failure_counter} times this epoch, epoch DONE :)")
        
main()