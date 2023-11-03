import pandas as pd

# Load the CSV data into a DataFrame
csv_file_path = '/home/lsalab/gpupool-with-k8s/experiment/reconfig_log.csv'
df = pd.read_csv(csv_file_path)

# 1. Remove rows with the same timestamp
df = df.drop_duplicates(subset='timestamp')

# 2. Convert all duration values to seconds
def convert_to_seconds(duration_str):
    if 'ms' in duration_str:
        return float(duration_str.replace('ms', '')) / 1000  # Convert milliseconds to seconds
    else:
        return float(duration_str.replace('s', ''))

df['duration'] = df['duration'].apply(convert_to_seconds)

# 3. Calculate the average duration for each event type
average_durations = df.groupby('event_type')['duration'].mean().reset_index()

# Print the DataFrame with average durations
print(average_durations)
