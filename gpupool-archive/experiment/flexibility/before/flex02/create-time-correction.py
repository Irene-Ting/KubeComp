import pandas as pd

offset_df = pd.read_csv('/home/lsalab/gpupool-with-k8s/experiment/flexibility/flex02/job_log_with_offset.csv')
pod_in_fix_cluster = pd.read_csv('/home/lsalab/gpupool-with-k8s/experiment/flexibility/flex02/pod-cluster-4-4.csv')

# Merge the DataFrames on the 'Job' column
combined_df = pd.merge(offset_df, pod_in_fix_cluster, on='Job', how='inner')
columns_to_keep = ['Job', 'Offset', 'Create', 'Start', 'Wait']
combined_df = combined_df[columns_to_keep]

combined_df['Create'] = pd.to_datetime(combined_df['Create'])
combined_df['Start'] = pd.to_datetime(combined_df['Start'])
combined_df['Offset'] = pd.to_timedelta(combined_df['Offset'], unit='s')

start = combined_df['Create'][0]
combined_df['Real Wait'] = combined_df['Start'] - (start + combined_df['Offset'])
combined_df['Real Wait'] = pd.to_timedelta(combined_df['Real Wait'], unit='s').dt.total_seconds()

# If needed, save the combined DataFrame to a new CSV file
combined_df.to_csv('/home/lsalab/gpupool-with-k8s/experiment/flexibility/flex02/pod-cluster-4-4-combined.csv', index=False)

# Print the combined DataFrame
# print(combined_df)
