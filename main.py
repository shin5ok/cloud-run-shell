import os
from fastapi import FastAPI, Depends, Header, Request, Response, APIRouter, status
from fastapi.responses import JSONResponse
from pydantic import *
from typing import Union, List
import uvicorn
import subprocess
from subprocess import PIPE
import requests

app = FastAPI()
secret = os.environ.get("SECRET", "gcp")
port = os.environ.get("PORT", "8080")
secret_header_key = "X-MyGCP-Secret"
my_id: str = ""

def _get_my_id() -> str:
    global my_id
    if my_id:
        return my_id
    try:
        headers = {'Metadata-Flavor':'Google'}
        response = requests.get("http://metadata.google.internal/computeMetadata/v1/instance/id", headers=headers, timeout=0.1)
        my_id = response.content
        return my_id
    except:
        return "__no_name__"

class Command(BaseModel):
    secret: Union[str, None] = None
    command: str

class Output(BaseModel):
    stdout_output: List[str] = []
    stderr_output: List[str] = []
    message: str = ""
    return_code: int = 0
    metadata: List[dict[str, str]]

@app.post("/shellcommand")
def cmd(command: Command, request: Request, response: Response, x_mygcp_secret = Header(default=None)):

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

@app.get("/longlong")
@app.get("/longlong/{s}")
def long(s: int = 1):
    import time
    time.sleep(s)
    return {"wait":s}

@app.middleware("http")
async def simple_auth(request: Request, call_next):
    response = await call_next(request)
    secret_in_request = request.headers.get(secret_header_key)

    if not secret_in_request:
        message = "auth required"
        return JSONResponse(
            dict(message=message),
            status.HTTP_401_UNAUTHORIZED
        )

    if secret != secret_in_request:
        message = "auth error"
        return JSONResponse(
            dict(message=message),
            status.HTTP_403_FORBIDDEN
        )

    return response

if __name__ == "__main__":
    options = {
            'port': int(port),
            'host': '0.0.0.0',
            'workers': 256,
            'reload': True,
        }
    uvicorn.run("main:app", **options)
