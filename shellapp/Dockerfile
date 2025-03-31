FROM python:3.12-slim

WORKDIR /shell
COPY . .

RUN pip install --no-cache-dir poetry==1.8.5; \
        poetry config virtualenvs.in-project true


RUN poetry install; \
    apt-get update && \
    apt-get install -y curl iproute2 git procps net-tools build-essential iperf3 qperf wget apache2 golang sysbench imagemagick ffmpeg ; \
    apt-get install -y lsb-release time nfs-common libreadline-dev zlib1g-dev; \
    gcsFuseRepo=gcsfuse-`lsb_release -c -s`; \
    echo "deb https://packages.cloud.google.com/apt $gcsFuseRepo main" | \
    tee /etc/apt/sources.list.d/gcsfuse.list; \
    curl https://packages.cloud.google.com/apt/doc/apt-key.gpg | \
    apt-key add -; \
    apt-get update; \
    apt-get install -y gcsfuse && \
    apt-get clean

# RUN git clone https://github.com/kdlucas/byte-unixbench && cd byte-unixbench/UnixBench && make

ENV CLOUD_SHELL_SRC https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-cli-457.0.0-linux-x86_64.tar.gz
RUN curl $CLOUD_SHELL_SRC | tar xzf - && ./google-cloud-sdk/install.sh ; \
    ./google-cloud-sdk/bin/gcloud init; true

ENV PATH /shell/google-cloud-sdk/bin:$PATH

ENV PYTHONUNBUFFERED true

CMD ["poetry", "run", "python", "main.py"]
