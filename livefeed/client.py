""" This is a module for processing messages """
import json
import gdax
import requests

with open('secret.json') as data_file:
    DATA = json.load(data_file)

class Client(gdax.WebsocketClient):
    """
    Class: Client - overrides default methods to provide
    necessary implementations
    """
    def __init__(self, url="wss://ws-feed.gdax.com", products=None, message_type="subscribe"):
        super(Client, self).__init__(url, products, message_type)

    def on_message(self, msg):
        if msg["type"] == "match":
            value = float(msg["price"])
            if value_is_abnormal(value):
                requests.post(DATA["url"]+DATA["key"], params={'value1' : str(value)})

def value_is_abnormal(value):
    """
    Determines if the given value is worthy of notifying the user
    """
    return True