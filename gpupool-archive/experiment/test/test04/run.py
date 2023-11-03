import numpy as np
import glob 
import time
import os

cmd = f"kubectl create -f ./input/pod-06.yaml"
print(cmd)
os.system(cmd)
time.sleep(10)

cmd = f"kubectl create -f ./input/pod-07.yaml"
print(cmd)
os.system(cmd)
time.sleep(10)

cmd = f"kubectl create -f ./input/pod-08.yaml"
print(cmd)
os.system(cmd)
time.sleep(10)

cmd = f"kubectl create -f ./input/pod-09.yaml"
print(cmd)
os.system(cmd)