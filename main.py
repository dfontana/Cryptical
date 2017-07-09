""" Main.py """
from client import Client

WSCLIENT = Client(url="wss://ws-feed.gdax.com", products="ETH-USD")
WSCLIENT.start()
