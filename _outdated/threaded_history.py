""" This is a script for obtaining time series data from GDAX """
import datetime
import os
import shutil
import csv
from multiprocessing.pool import ThreadPool
import gdax

CLIENT = gdax.PublicClient()
GRANULARITY = 1 # second
WAVE_SIZE = 7

def main(srttime=None, endtime=None):
    """
    Breaks down the given time period into digestable request "chunks" that
    the GDAX API can process. Outputs results into a CSV file.
    """
    with open('part_master.csv', 'w') as the_file:
        writer = csv.writer(the_file, dialect='excel')
        writer.writerow(['time', 'low', 'high', 'open', 'close', 'volume'])

    requests = (endtime-srttime).total_seconds() / GRANULARITY

    # Build thread queue
    print("Constructing Ship...")
    threads = define_threads(requests, srttime, endtime)
    if not os.path.exists("parts"):
        os.makedirs("parts")
    out_file = open("part_master.csv", "a")

    # Unleash the threads
    print("Unleashing the Kraken...")
    pool = ThreadPool(processes=4)
    pool.map(lambda s: process_frame(s[0], s[1], s[2]), threads)
    pool.close()
    pool.join()

    # Seal the deal.
    write_parts_to_master(out_file)
    out_file.close()
    print("The Seas Have Settled.")


def define_threads(requests, srttime, endtime):
    """
    Builds an array of threads to process, where each thread handles a chunk of time
    """
    ths = []
    count = 0
    if requests > 200:
        sframe = srttime
        eframe = sframe + datetime.timedelta(seconds=GRANULARITY*200)
        while eframe <= endtime:
            ths.append([sframe, eframe, count])
            sframe = eframe + datetime.timedelta(seconds=GRANULARITY)
            eframe = sframe + datetime.timedelta(seconds=GRANULARITY*200)
            count += 1
        if eframe > endtime:
            ths.append([sframe, eframe, count])
    else:
        ths.append([sframe, eframe, count])
    return ths



def write_parts_to_master(out_file):
    """
    Writes the current contents of the parts directory to the master csv.
    After which, it deletes the parts & rebuilds folder for the next wave.
    """
    for filename in os.listdir("parts"):
        with open("parts/"+filename) as part:
            for line in part:
                out_file.write(line)



def process_frame(start_frame, end_frame, thread_count):
    """
    Makes a call to the historic endpoint for the given time period, writing results
    to file. Sometimes the API returns "message" - that data row is filtered out.
    Additionally, the timestamp is in epoch time - which has been converted to
    human readable output in UTC time.
    """
    print(start_frame, end_frame)
    subarray = CLIENT.get_product_historic_rates('ETH-USD', start=start_frame,
                                                 end=end_frame, granularity=GRANULARITY)

    with open('parts/part_'+str(thread_count)+'.csv', 'w') as the_file:
        writer = csv.writer(the_file, dialect='excel')
        for row in subarray:
            if row[0] == 'm':
                break
            row[0] = datetime.datetime.fromtimestamp(row[0]).strftime('%x %X')
            writer.writerow(row)

START = datetime.datetime(2017, 7, 30, hour=19)
END = datetime.datetime.now()
main(START, END)
