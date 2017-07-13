""" This is a module for processing messages """
import datetime
import gdax

CLIENT = gdax.PublicClient()

def main(srttime=None, endtime=None):
    """
    Doc string
    """
    gran = 1 # second

    result_dict = {}

    requests = (endtime-srttime).total_seconds() / gran
    if requests > 200:
        start_frame = srttime
        end_frame = start_frame + datetime.timedelta(seconds=gran*200)
        while end_frame <= endtime:
            result_dict.update(CLIENT.get_product_historic_rates('ETH-USD', start=start_frame, end=end_frame, granularity=gran))
            start_frame = end_frame + datetime.timedelta(seconds=gran)
            end_frame = start_frame + datetime.timedelta(seconds=gran*200)
        if end_frame > endtime:
            # there's extra, need one last request
            result_dict.update(CLIENT.get_product_historic_rates('ETH-USD', start=start_frame, end=endtime, granularity=gran))
    else:
        result_dict.update(CLIENT.get_product_historic_rates('ETH-USD', start=srttime, end=endtime, granularity=gran))
    return result_dict

START = datetime.datetime(2017, 1, 1, 8, 0)
END = datetime.datetime.now()
print main(START, END)
