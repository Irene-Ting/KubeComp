import csv
import sys

# Existing CSV file path
csv_file_path = '/home/lsalab/gpupool-with-k8s/experiment/scheduler_log.csv'

# Initialize an empty list to store the new log entries as dictionaries
new_log_entries = []

log_file = '/home/lsalab/gpupool-with-k8s/experiment/scheduler.log'
# Open the log file for reading
with open(log_file, 'r') as file:
    for line in file:
        parts = line.strip().split()
        if parts[0] != 'I1016':
            continue
        if (len(parts) > 6 and parts[6] == "Prefilter") or (len(parts) > 3 and parts[3] == "event.go:307]"):
            entry = {
                'timestamp': parts[1],
                'description': ' '.join(parts[3:]),
            }
            new_log_entries.append(entry)

# Append the new data to the existing CSV file
with open(csv_file_path, 'a', newline='') as csv_file:
    fieldnames = ['timestamp', 'description']
    writer = csv.DictWriter(csv_file, fieldnames=fieldnames)

    for entry in new_log_entries:
        writer.writerow(entry)

print(f'New data has been appended to {csv_file_path}')
