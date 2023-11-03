FROM python:3.10-alpine

WORKDIR /app

COPY requirements.txt ./
RUN pip install --no-cache-dir -r requirements.txt

COPY gpu_attached_exporter.py /app

EXPOSE 8001

CMD [ "python", "gpu_attached_exporter.py" ]

# docker build . -t gpu_attached_exporter
# docker run -it --name "gpu_attached_exporter" -p 127.0.0.1:8001:8001 -e DEBUG=1  --rm gpu_attached_exporter:latest
# docker kill --signal=SIGKILL "gpu_attached_exporter"