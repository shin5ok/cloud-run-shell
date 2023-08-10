FROM python:3.10-slim

WORKDIR /shell
COPY . .
RUN pip install --no-cache-dir poetry
RUN poetry config virtualenvs.in-project true
RUN poetry install; \
    apt-get update && apt-get install -y curl iproute2 git procps net-tools build-essential iperf3 qperf wget

RUN git clone https://github.com/kdlucas/byte-unixbench && cd byte-unixbench/UnixBench && make && cp Run /usr/local/bin/Run

CMD ["poetry", "run", "python", "main.py"]
