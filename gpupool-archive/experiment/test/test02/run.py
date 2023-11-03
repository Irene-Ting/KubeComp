import numpy as np
import glob 
import time
import os

def run_experiment(folder, lam):
  files = glob.glob(f"{folder}/*")
  files.sort()
  for f in files:
    cmd = f"kubectl create -f {f}"
    # cmd = "ls"
    print(cmd)
    os.system(cmd)
    time.sleep(np.random.poisson(lam=lam))

folder = "/home/lsalab/gpupool-with-k8s/experiment/test02/input"
lam = 30
run_experiment(folder, lam)