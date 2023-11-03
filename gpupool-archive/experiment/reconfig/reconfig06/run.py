import numpy as np
import glob 
import time
import os

for i in range(3):
    cmd = f"kubectl create -f ./input/pod-01.yaml"
    print(cmd)
    os.system(cmd)
    time.sleep(250)

    cmd = f"kubectl create -f ./input/pod-03.yaml"
    print(cmd)
    os.system(cmd)
    time.sleep(250)

    cmd = f"kubectl delete -f ./input"
    print(cmd)
    os.system(cmd)
    time.sleep(10)

# cmd = f"kubectl create -f ./input/pod-03.yaml"
# print(cmd)
# os.system(cmd)
# time.sleep(30)

# cmd = f"kubectl create -f ./input/pod-04.yaml"
# print(cmd)
# os.system(cmd)
# time.sleep(60)

# cmd = f"kubectl create -f ./input/pod-05.yaml"
# print(cmd)
# os.system(cmd)