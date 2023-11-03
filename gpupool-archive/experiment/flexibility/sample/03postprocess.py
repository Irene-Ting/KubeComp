import json
import csv
import pandas as pd
import sys
import os

file_name = f"pod-cluster-{sys.argv[1]}"
json_file = f"./{file_name}.json"

cmd = f"kubectl get po -o json > {json_file}"
print(cmd)
os.system(cmd)

with open(json_file, "r") as file:
    json_data = file.read()
data = json.loads(json_data)

df = pd.DataFrame(columns=["Job", "Create", "Start"])

for i in range(len(data["items"])):
  creation_timestamp = pd.to_datetime(data["items"][i]["metadata"]["creationTimestamp"])
  name = data["items"][i]["metadata"]["name"][-3:]
  start_time = pd.to_datetime(data["items"][i]["status"]["startTime"])
  df.loc[len(df)] = [name, creation_timestamp.strftime("%Y-%m-%d %H:%M:%S"), start_time.strftime("%Y-%m-%d %H:%M:%S")]

df.to_csv(f"./{file_name}.csv")

# create time correction
offset_df = pd.read_csv('./job_log.csv')
pod_in_fix_cluster = pd.read_csv(f"./{file_name}.csv")

# Merge the DataFrames on the 'Job' column
combined_df = pd.merge(offset_df, pod_in_fix_cluster, on='Job', how='inner')
columns_to_keep = ['Job', 'Create', 'Start', 'Offset', 'Sleep']
combined_df = combined_df[columns_to_keep]

combined_df['Create'] = pd.to_datetime(combined_df['Create'])
combined_df['Start'] = pd.to_datetime(combined_df['Start'])
combined_df['Offset'] = pd.to_timedelta(combined_df['Offset'], unit='s')

start = combined_df['Create'][0]
combined_df['Real Start'] = start + combined_df['Offset']
combined_df['Real Wait'] = combined_df['Start'] - combined_df['Real Start']
combined_df['Real Wait'] = pd.to_timedelta(combined_df['Real Wait'], unit='s').dt.total_seconds()

# If needed, save the combined DataFrame to a new CSV file
combined_df.to_csv(f"./{file_name}-real-wait.csv", index=False)
