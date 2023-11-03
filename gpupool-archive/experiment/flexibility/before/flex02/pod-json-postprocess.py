import json
import csv
import pandas as pd
file_name = "pod-cluster-2-6"
json_file = f"/home/lsalab/gpupool-with-k8s/experiment/flexibility/flex02/{file_name}.json"
with open(json_file, "r") as file:
    json_data = file.read()
data = json.loads(json_data)

df = pd.DataFrame(columns=["Job", "Create", "Start", "Wait"])

for i in range(len(data["items"])):
  creation_timestamp = pd.to_datetime(data["items"][i]["metadata"]["creationTimestamp"])
  name = data["items"][i]["metadata"]["name"]
  start_time = pd.to_datetime(data["items"][i]["status"]["startTime"])
  wait_time = (start_time - creation_timestamp).total_seconds()
  df.loc[len(df)] = [name, creation_timestamp, start_time.strftime("%Y-%m-%d %H:%M:%S"), wait_time]

df.to_csv(f"/home/lsalab/gpupool-with-k8s/experiment/flexibility/flex02/{file_name}.csv")