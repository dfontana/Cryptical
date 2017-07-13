""" This is a script for obtaining time series data from GDAX """
import datetime
import csv
import gdax

CLIENT = gdax.PublicClient()
GRANULARITY = 1 # second
OUT_FILE = open("data_points.csv", 'w')
WRITER = csv.writer(OUT_FILE, dialect='excel')

def main(srttime=None, endtime=None):
    """
    Breaks down the given time period into digestable request "chunks" that
    the GDAX API can process. Outputs results into a CSV file.
    """
    WRITER.writerow(['time', 'low', 'high', 'open', 'close', 'volume'])

    requests = (endtime-srttime).total_seconds() / GRANULARITY

    if requests > 200:
        start_frame = srttime
        end_frame = start_frame + datetime.timedelta(seconds=GRANULARITY*200)
        while end_frame <= endtime:
            process_time_frame(start_frame, end_frame)
            start_frame = end_frame + datetime.timedelta(seconds=GRANULARITY)
            end_frame = start_frame + datetime.timedelta(seconds=GRANULARITY*200)
        if end_frame > endtime:
            process_time_frame(start_frame, endtime)
    else:
        process_time_frame(srttime, endtime)
    OUT_FILE.close()

def process_time_frame(start_frame, end_frame):
    """
    Makes a call to the historic endpoint for the given time period, writing results
    to file. Sometimes the API returns "message" - that data row is filtered out.
    Additionally, the timestamp is in epoch time - which has been converted to
    human readable output in UTC time.
    """
    subarray = CLIENT.get_product_historic_rates('ETH-USD', start=start_frame,
                                                 end=end_frame, granularity=GRANULARITY)
    for row in subarray:
        if row[0] == 'm':
            break
        row[0] = datetime.datetime.fromtimestamp(row[0]).strftime('%x %X')
        WRITER.writerow(row)

START = datetime.datetime(2017, 7, 12, 21, 0)
END = datetime.datetime(2017, 7, 12, 22, 0)
main(START, END)
