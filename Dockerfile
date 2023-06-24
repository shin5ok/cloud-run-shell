FROM python:3.10-slim

COPY . .
RUN pip install poetry
RUN poetry install
RUN apt update && apt install -y curl iproute2 git procps net-tools build-essential

RUN git clone https://github.com/kdlucas/byte-unixbench && cd byte-unixbench/UnixBench && make && cp Run /usr/local/bin/Run

CMD ["poetry", "run", "python", "main.py"]
