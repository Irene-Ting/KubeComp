import csv
import sys

# Existing CSV file path
csv_file_path = '/home/lsalab/gpupool-with-k8s/experiment/reconfig_log.csv'

# Initialize an empty list to store the new log entries as dictionaries
new_log_entries = []

# Open the log file for reading
with open(sys.argv[1], 'r') as file:
    for line in file:
        parts = line.strip().split()
        if parts[-1] == 'Reconfig':
            continue
        entry = {
            'timestamp': parts[0] + ' ' + parts[1],
            'event_type': ' '.join(parts[2:-1]),
            'duration': parts[-1]
        }
        new_log_entries.append(entry)

# Append the new data to the existing CSV file
with open(csv_file_path, 'a', newline='') as csv_file:
    fieldnames = ['timestamp', 'event_type', 'duration']
    writer = csv.DictWriter(csv_file, fieldnames=fieldnames)

    for entry in new_log_entries:
        writer.writerow(entry)

print(f'New data has been appended to {csv_file_path}')
