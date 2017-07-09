""" This is a module for processing messages """
import gdax

class Client(gdax.WebsocketClient):
    """
    Class: Client - overrides default methods to provide
    necessary implementations
    """
    def __init__(self, url="wss://ws-feed.gdax.com", products=None, message_type="subscribe"):
        self.message_count = 0
        super(Client, self).__init__(url, products, message_type)

    def on_message(self, msg):
        if msg["type"] == "match":
            value = float(msg["price"])
            print("Value: ", value)
            self.message_count += 1
