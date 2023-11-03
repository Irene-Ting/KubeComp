import threading
import queue
import pandas as pd
import time
import os
import subprocess
import numpy as np
import csv
from datetime import datetime
import sys

WORKLOADDIR = "/home/lsalab/gpupool-with-k8s/experiment/flexibility/workload"
csv_file = "./job_log.csv"
df = pd.read_csv(csv_file)

def log_cluster_event(job):
    timestamp = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
    with open(f"./run-gpu-pool.csv", mode='a', newline='') as file:
        writer = csv.writer(file)
        writer.writerow([timestamp, job])

# Function for the first thread (producer)
def producer():
    for index, row in df.iterrows():
        job = row['Job']
        sleep = row['Sleep']
        time.sleep(int(float(sleep)))

        cmd = f'kubectl create -f {WORKLOADDIR}/pod-{job}.yaml'
        buffer.put(cmd)
        log_cluster_event(job)
    

def consumer():
    while True:
        try:
            # Get an item from the buffer
            cmd = buffer.get(timeout=120)  # Adjust the timeout as needed
            while True:
                output = subprocess.check_output("kubectl get po", shell=True, text=True)
                pods_info = output.strip()
                if any(status in pods_info for status in ("Pending", 'ContainerCreating')):
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
