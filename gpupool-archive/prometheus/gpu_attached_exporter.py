from prometheus_client import start_http_server, Gauge
import time
import os
import subprocess
import yaml

def get_gpu():
    cmd = "kubectl get node -o yaml"
    output_bytes = subprocess.check_output(cmd, shell=True)
    output_yaml = output_bytes.decode("utf-8")

    data = yaml.safe_load(output_yaml)
    for node in data.get('items', []):
        node_name = node.get('metadata', {}).get('name', '')
        if node_name == "css-host-128":
            node1 = node.get('status', '').get('allocatable', '').get('falcon.com/gpu', '')
        elif node_name == "css-host-129":
            node2 = node.get('status', '').get('allocatable', '').get('falcon.com/gpu', '')

    return node1, node2


# def get_gpu():
#     node1 = 2
#     node2 = 6
#     # cmd = "ls"
#     print(cmd)
#     os.system(cmd)

#     return node1, node2

gpu_count = Gauge('node_gpu_count', 'Number of GPUs on the node', ['server_name'])

if __name__ == '__main__':

    # get_gpu_test()

    start_http_server(8001)

    while True:
        v1, v2 = get_gpu()
        
        # Set the GPU count metric with the server name as a label
        gpu_count.labels(server_name="css-host-128").set(v1)
        gpu_count.labels(server_name="css-host-129").set(v2)

        # Sleep for a while (adjust the sleep duration as needed)
        time.sleep(1)
