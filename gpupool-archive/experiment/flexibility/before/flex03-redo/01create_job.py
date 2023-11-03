import os
import time
import glob
import numpy as np
import random
import csv
from datetime import datetime
import pandas as pd

WORKLOADDIR = "/home/lsalab/gpupool-with-k8s/experiment/flexibility/workload"
csv_file = "./job_log.csv"

def log_cluster_event(job):
    timestamp = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
    with open(csv_file, mode='a', newline='') as file:
        writer = csv.writer(file)
        writer.writerow([timestamp, job])

def gen_workload(total, percentage):
    work_list = []
    for i in range(len(percentage)):
        work_cnt = int(total * percentage[i])
        for j in range(1, work_cnt + 1):
            work_list.append(chr(97+i)+str(j).zfill(2))
    random.shuffle(work_list)
    return work_list

workload = gen_workload(20, [0.25, 0.25, 0.25, 0.25, 0])

df = pd.DataFrame(columns=["Job", "Offset", "Sleep"])

wait = 0
total_time = 0
for w in workload:
    offset = total_time
    next_wait = np.random.poisson(lam=60)
    df.loc[len(df)] = [w, offset, wait]

    wait = next_wait
    total_time += wait

df.to_csv(csv_file, index=False)
