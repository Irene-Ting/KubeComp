from prometheus_client import start_http_server, Gauge
import subprocess
import os
import time
import sys

# IPMIUTIL_COMMAND=os.environ["IPMIUTIL_COMMAND"]
# print("Command:", IPMIUTIL_COMMAND, file=sys.stderr)

DEBUG="DEBUG" in os.environ and os.environ["DEBUG"] == "1"

gpu_cnt = Gauge('gpu_cnt', 'get GPU cnt from Falcon API')

def get_gpus():
    # output = subprocess.check_output(['/bin/sh', '-c', IPMIUTIL_COMMAND]).strip()
    output = 3
    if DEBUG:
        print(output, file=sys.stderr)
    return int(output)

gpu_cnt.set_function(get_gpus)

if __name__ == '__main__':
    start_http_server(8001)

    while True:
        time.sleep(1e9)