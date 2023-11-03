import numpy as np
import glob 
import time
import os

cmd = f"kubectl create -f ./input/pod-01.yaml"
print(cmd)
os.system(cmd)
time.sleep(90)

cmd = f"kubectl create -f ./input/pod-02.yaml"
print(cmd)
os.system(cmd)
time.sleep(5)

cmd = f"kubectl create -f ./input/pod-03.yaml"
print(cmd)
os.system(cmd)