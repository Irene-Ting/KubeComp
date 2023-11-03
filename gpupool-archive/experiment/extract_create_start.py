import json
import csv
import pandas as pd
import sys
import os

file_name = sys.argv[1]
json_file = f"./{file_name}.json"

cmd = f"kubectl get po -o json > {json_file}"
print(cmd)
os.system(cmd)

with open(json_file, "r") as file:
    json_data = file.read()
data = json.loads(json_data)

df = pd.DataFrame(columns=["Job", "Create", "Start", "Wait"])

for i in range(len(data["items"])):
  creation_timestamp = pd.to_datetime(data["items"][i]["metadata"]["creationTimestamp"])
  name = data["items"][i]["metadata"]["name"]
  start_time = pd.to_datetime(data["items"][i]["status"]["startTime"])
  wait = (start_time - creation_timestamp).total_seconds()
  df.loc[len(df)] = [name, creation_timestamp.strftime("%Y-%m-%d %H:%M:%S"), start_time.strftime("%Y-%m-%d %H:%M:%S"), wait]

df.to_csv(f"./{file_name}.csv", index=False)
