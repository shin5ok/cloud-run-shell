import os
from fastapi import FastAPI, Header, Request, Response, status
from fastapi.responses import JSONResponse
from pydantic import *
import subprocess
from subprocess import PIPE
import requests, logging, sys

secret = os.environ.get("SECRET")
listen_port = os.environ.get("PORT", "8080")
secret_header_key = "X-MyGCP-Secret"

logger = logging.getLogger("uvicorn")
handler = logging.StreamHandler(sys.stdout)
logger.addHandler(handler)
logger.setLevel(logging.INFO)

app = FastAPI()

my_id: str = ""

def _get_my_id() -> str:
    global my_id
    if my_id:
        return my_id
    try:
        headers = {'Metadata-Flavor':'Google'}
        response = requests.get("http://metadata.google.internal/computeMetadata/v1/instance/id", headers=headers, timeout=0.1)
        return response.content
    except:
        return "__no_name__"

class Command(BaseModel):
    secret: str | None = None
    command: str

class Output(BaseModel):
    stdout_output: list[str] = []
    stderr_output: list[str] = []
    message: str = ""
    return_code: int = 0
    metadata: list[dict[str, str]]

@app.post("/shellcommand")
def cmd(command: Command, request: Request, response: Response, x_mygcp_secret = Header(default=None)):

    logger.info(command)

    cmd_str = command.command

    proc = subprocess.run(cmd_str, shell=True, stdout=PIPE, stderr=PIPE, text=True)
    stdout_output = proc.stdout.split("\n")
    stderr_output = proc.stderr.split("\n")

    return Output(
            stderr_output=stderr_output,
            stdout_output=stdout_output,
            return_code=proc.returncode,
            metadata=[dict(instance_id=_get_my_id())],
        )

@app.get("/ping")
def _ping():
    logger.info("GET /ping")
    return JSONResponse(
        "pong",
        status.HTTP_200_OK
    )

@app.middleware("http")
async def simple_auth(request: Request, call_next):
    response = await call_next(request)

    if secret:
        if not secret_in_request:
            message = "auth required"
            return JSONResponse(
                dict(message=message),
                status.HTTP_401_UNAUTHORIZED
            )

        secret_in_request = request.headers.get(secret_header_key)
        if secret != secret_in_request:
            message = "auth error"
            return JSONResponse(
                dict(message=message),
                status.HTTP_403_FORBIDDEN
            )

    return response

def main():
    import asyncio

    from hypercorn.config import Config
    from hypercorn.asyncio import serve

    config = Config()
    config.bind = [f"0.0.0.0:{listen_port}"]
    config.keep_alive_timeout = 3600
    config.read_timeout = 3600
    asyncio.run(serve(app, config))

if __name__ == "__main__":
    main()
