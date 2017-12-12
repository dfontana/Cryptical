""" This is a module for processing messages """
import json
import gdax
import requests

with open('secret.json') as data_file:
    """
    Used by IFTTT service to send notification to webhook of value change
    thus notifying the user. This file would provide the URL and user key
    of the IFTTT service.
    """
    DATA = json.load(data_file)

class Client(gdax.WebsocketClient):
    """
    Class: Client - overrides default methods to provide
    necessary implementations
    """
    def __init__(self, url="wss://ws-feed.gdax.com", products=None, message_type="subscribe"):
        super(Client, self).__init__(url, products, message_type)

    def on_message(self, msg):
        """
        Handles messages arriving to client from socket stream, we
        only listen to "match" messages - meaning a buy/sell was successful
        """
        if msg["type"] == "match":
            value = float(msg["price"])

            # Notify user only when they should know
            if value_is_abnormal(value):
                requests.post(DATA["url"]+DATA["key"], params={'value1' : str(value)})
            
            # Consider value with your rolling average
            add_to_average(value)

def value_is_abnormal(value):
    """
    Determines if the given value stands apart from the running average.
    We don't add the value to the average before doing this, as we don't 
    want to skew the average right before comparing against it.

    TODO Needs to be implemented - this would be the signaller that current
    prices are out of average.
    """
    return True

def add_to_average(value):
    """
    Takes the given value and adds it into the running average, thus updating
    what the the script knows currency to be worth.

    TODO Needs to be implemented. Consider https://www.investopedia.com/terms/b/bollingerbands.asp
    """
    return None