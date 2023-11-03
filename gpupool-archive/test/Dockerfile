FROM nvidia/cuda:11.1.1-devel-ubuntu20.04
COPY run /bin/gpu-run
CMD /bin/gpu-run 4 5
# docker run --rm gpu-test /bin/bash -c "nvidia-smi  && sleep 3"