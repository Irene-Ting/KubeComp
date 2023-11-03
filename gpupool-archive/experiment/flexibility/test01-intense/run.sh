python3 02-2-run-fix.py 4-4
sleep 2000
python3 03postprocess.py 4-4
kubectl delete -f ../workload
sleep 20
python3 02-1-run-gpu-pool.py
sleep 2000
python3 03postprocess.py pool