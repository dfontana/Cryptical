""" This is a script for obtaining time series data from GDAX """
import datetime
import os
import shutil
import csv
import threading
from itertools import zip_longest
import gdax


"""

TODO:
    - Make the parts folder if its not on disk
    - Concatenate all files of a wave into the HEADER file after they finish
    and before starting the next wave.
        - Consider doing it every X number of waves, wiping the parts folder after
        the merge, and then when all threads are done processing in full concatenate
        any remaining files in the parts folder before deleting it.

"""


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
    threads = []

    # Build thread queue
    print("Constructing Threads...")
    thread_count = 0
    if requests > 200:
        start_frame = srttime
        end_frame = start_frame + datetime.timedelta(seconds=GRANULARITY*200)
        while end_frame <= endtime:
            thread = threading.Thread(target=process_time_frame, args=(start_frame, end_frame, thread_count))
            threads.append(thread)
            thread_count += 1
            start_frame = end_frame + datetime.timedelta(seconds=GRANULARITY)
            end_frame = start_frame + datetime.timedelta(seconds=GRANULARITY*200)
        if end_frame > endtime:
            thread = threading.Thread(target=process_time_frame, args=(start_frame, endtime, thread_count))
            threads.append(thread)
            thread_count += 1
    else:
        thread = threading.Thread(target=process_time_frame, args=(srttime, endtime, thread_count))
        threads.append(thread)
        thread_count += 1

    # Unleash the threads
    print(str(len(threads)) + " threads constructed.")
    print("Unleashing the Kraken (In waves)...")
    
    if not os.path.exists("parts"):
        os.makedirs("parts")
    out_file = open("part_master.csv", "a")

    wave_index = 1
    wave_size = 0
    for group in grouper(WAVE_SIZE, threads):
        print("\tStarting Wave " + str(wave_index) + "/" + str(len(threads)))
        for thr in group:
            if thr is None:
                continue
            else:
                wave_size += 1
                thr.start()

        for thr in group:
            if thr is None:
                continue
            else:
                thr.join()

        print("\tStitching the masts")
        write_parts_to_master(out_file)
        wave_index += 1

    # Seal the deal.
    out_file.close()
    print("The Seas Have Settled.")

def write_parts_to_master(out_file):
    """
    Writes the current contents of the parts directory to the master csv.
    After which, it deletes the parts & rebuilds folder for the next wave.
    """
    for filename in os.listdir("parts"):
        with open("parts/"+filename) as part:
            for line in part:
                out_file.write(line)
    shutil.rmtree("parts")
    os.makedirs("parts")

def process_time_frame(start_frame, end_frame, thread_count):
    """
    Makes a call to the historic endpoint for the given time period, writing results
    to file. Sometimes the API returns "message" - that data row is filtered out.
    Additionally, the timestamp is in epoch time - which has been converted to
    human readable output in UTC time.
    """
    subarray = CLIENT.get_product_historic_rates('ETH-USD', start=start_frame,
                                                 end=end_frame, granularity=GRANULARITY)

    with open('parts/part_'+str(thread_count)+'.csv', 'w') as the_file:
        writer = csv.writer(the_file, dialect='excel')
        for row in subarray:
            if row[0] == 'm':
                break
            row[0] = datetime.datetime.fromtimestamp(row[0]).strftime('%x %X')
            writer.writerow(row)


def grouper(chunk_size, iterable, fillvalue=None):
    """
    Splits an array into chunk_sized subarrays, filling in empty spaces
    with None by default.
    """
    args = [iter(iterable)] * chunk_size
    return zip_longest(fillvalue=fillvalue, *args)

START = datetime.datetime(2017, 1, 1, 6, 0)
END = datetime.datetime.now()
main(START, END)
