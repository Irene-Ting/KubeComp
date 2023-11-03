import pandas as pd

# Load the CSV file into a DataFrame
csv_file = "/home/lsalab/gpupool-with-k8s/experiment/flexibility/flex02/job_log.csv"  # Replace with the path to your CSV file
df = pd.read_csv(csv_file, header=None, names=['Timestamp', 'Job'])

# Convert the "Timestamp" column to datetime objects
df['Timestamp'] = pd.to_datetime(df['Timestamp'])

# Calculate the time difference and shift by one to get the sleep time
df['Offset'] = (df['Timestamp'] - df['Timestamp'][0]).dt.total_seconds()
df['Wait'] = (df['Timestamp'] - df['Timestamp'].shift(1)).dt.total_seconds()

# Save the DataFrame back to a new CSV file with the "Sleep Time" column
output_csv = "job_log_with_offset.csv"  # Replace with your desired output file name
df.to_csv(output_csv, index=False)

print(f"New CSV file with sleep time saved to '{output_csv}'")
