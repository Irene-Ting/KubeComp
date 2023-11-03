import json
import csv
import pandas as pd
import sys
import os

json_file = sys.argv[1]

with open(json_file, "r") as file:
    json_data = file.read()
data = json.loads(json_data)

df = pd.DataFrame(columns=["Job", "GPU", "Container_Start", "Container_End"])

for i in range(len(data["items"])):
    name = data["items"][i]["metadata"]["name"][-3:]
    gpu   = data["items"][i]["spec"]["containers"][0]["resources"]["requests"]["falcon.com/gpu"]
    start_time = pd.to_datetime(data["items"][i]["status"]["containerStatuses"][0]["state"]["terminated"]["startedAt"])
    end_time = pd.to_datetime(data["items"][i]["status"]["containerStatuses"][0]["state"]["terminated"]["finishedAt"])
    df.loc[len(df)] = [name, gpu, start_time.strftime("%Y-%m-%d %H:%M:%S"), end_time.strftime("%Y-%m-%d %H:%M:%S")]

df.to_csv(f"./test04-util-{json_file[:-5]}.csv", index=False)