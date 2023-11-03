import threading
import queue
import pandas as pd
import time
import os
import subprocess
import numpy as np
import csv
from datetime import datetime

WORKLOADDIR = "/home/lsalab/gpupool-with-k8s/experiment/flexibility/workload-no-falcon-schduler"
csv_file = "/home/lsalab/gpupool-with-k8s/experiment/flexibility/flex02/job_log_with_offset.csv"
df = pd.read_csv(csv_file)

def log_cluster_event(job):
    timestamp = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
    with open("/home/lsalab/gpupool-with-k8s/experiment/flexibility/flex02/run-fix-0-8.csv", mode='a', newline='') as file:
        writer = csv.writer(file)
        writer.writerow([timestamp, job])

# Function for the first thread (producer)
def producer():
    for index, row in df.iterrows():
        job = row['Job']
        wait = row['Wait']
        if pd.isna(wait):
            wait = 0
        time.sleep(int(float(wait)))
        cmd = f'kubectl create -f {WORKLOADDIR}/pod-{job[-3:]}.yaml'
        buffer.put(cmd)
        log_cluster_event(job)
        # print(cmd)
        # os.system(cmd)
    

def consumer():
    while True:
        try:
            # Get an item from the buffer
            cmd = buffer.get(timeout=120)  # Adjust the timeout as needed
            while True:
                output = subprocess.check_output("kubectl get po", shell=True, text=True)
                pods_info = output.strip()
                if "Pending" in pods_info:
                    time.sleep(1)
                else:
                    print(cmd)
                    os.system(cmd)
                    break

        except queue.Empty:
            # If the buffer is empty, break out of the loop
            break

# Create a buffer (queue) to share between threads
buffer = queue.Queue()

# Create and start the producer thread
producer_thread = threading.Thread(target=producer)
producer_thread.start()

# Create and start the consumer thread
consumer_thread = threading.Thread(target=consumer)
consumer_thread.start()

# Wait for both threads to finish
producer_thread.join()
consumer_thread.join()
