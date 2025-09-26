#!/usr/bin/env python3
import feedparser
import requests
import time
from datetime import datetime

def main():
    feeders = [
        'https://news.un.org/feed/subscribe/en/news/region/global/feed/rss.xml',
        'http://rss.cnn.com/rss/cnn_topstories.rss',
        'https://feeds.bbci.co.uk/news/rss.xml',
        'https://www.cbsnews.com/latest/rss/main',
        'https://feeds.skynews.com/feeds/rss/world.xml'
    ]

    newsList = []
    for feedUrl in feeders:
        feed = feedparser.parse(feedUrl)
        for news in feed.entries:
            dt = datetime.now()
            try:
                dt = datetime.strptime(news.published, "%a, %d %b %Y %H:%M:%S %z")
            except:
                try:
                    dt = datetime.strptime(news.published, "%a, %d %b %Y %H:%M:%S GMT")
                except:
                    continue

            newsList.append({
                "message": news.title + ": " + news.link + " (" + news.published + ")",
                "timestamp": dt.timestamp()
            })

    newsList.sort(key=lambda news: news['timestamp'], reverse=False)
    for news in newsList:
        requests.post('http://127.0.0.1/chats', data={'msg': news['message']})
        print(news['message'])

main()
