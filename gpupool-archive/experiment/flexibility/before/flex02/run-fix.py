import os
import time
import subprocess
import numpy as np
import random
import csv
from datetime import datetime
import pandas as pd

WORKLOADDIR = "/home/lsalab/gpupool-with-k8s/experiment/flexibility/workload-no-falcon-schduler"
csv_file = "/home/lsalab/gpupool-with-k8s/experiment/flexibility/flex02/job_log_with_offset.csv"

df = pd.read_csv(csv_file)
start_time = datetime.now()

for index, row in df.iterrows():
    job = row['Job']
    offset = row['Offset']
    
    while True:
        output = subprocess.check_output("kubectl get po", shell=True, text=True)
        pods_info = output.strip()
        if "Pending" in pods_info:
            time.sleep(1)
        else:
            break

    while (datetime.now() - start_time).total_seconds() < offset:
        # print((datetime.now() - start_time).total_seconds())
        time.sleep(1)
    
    cmd = f'kubectl create -f {WORKLOADDIR}/pod-{job[-3:]}.yaml'
    print(cmd)
    os.system(cmd)
