import yaml
import os

def gen_yaml(index, gpu_num, cpu_num, duration):
  yaml = f'''apiVersion: v1
kind: Pod
metadata:
  name: falcon-pod-{index}-default
spec:
  containers:
  - name: falcon-pod-{index}-default
    image: "gpu-test:test"
    imagePullPolicy: Never
    command: ['sh', '-c', 'echo "start" && /bin/gpu-run {gpu_num} {duration}']
    resources:
      requests:
        memory: "500Mi"
        cpu: "{cpu_num}"
        falcon.com/gpu: "{gpu_num}"
      limits:
        memory: "500Mi"
        cpu: "{cpu_num}"
        falcon.com/gpu: "{gpu_num}"
  restartPolicy: OnFailure
  '''

  return yaml


cmd = 'mkdir -p "/home/lsalab/gpupool-with-k8s/experiment/overhead/gpupool-schd"'
os.system(cmd)

# for i in range(1,9):
#   with open(f'/home/lsalab/gpupool-with-k8s/experiment/overhead/gpupool-schd/pod-{str(i).zfill(2)}.yaml', 'w') as f:
#     f.write(gen_yaml(f'{str(i).zfill(2)}', 1, 1, 10))

for i in range(1,9):
  with open(f'/home/lsalab/gpupool-with-k8s/experiment/overhead/default-schd/pod-{str(i).zfill(2)}.yaml', 'w') as f:
    f.write(gen_yaml(f'{str(i).zfill(2)}', 0, 1, 10))