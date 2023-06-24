import os
from fastapi import FastAPI, Depends, Header, Request, Response, APIRouter, status
from pydantic import *
from typing import Union, List
import uvicorn
import json
import subprocess
from subprocess import PIPE

app = FastAPI()
secret = os.environ.get("SECRET")
port = os.environ.get("PORT", "8080")
secret_header_key = "X-MyGCP-Secret"

class Command(BaseModel):
    secret: Union[str, None] = None
    command: str

class Output(BaseModel):
    stdout_output: List[str] = []
    stderr_output: List[str] = []
    message: str = ""
    return_code: int = 0

@app.post("/shellcommand")
def cmd(command: Command, request: Request, response: Response, x_mygcp_secret = Header(default=None)):

    cmd_str = command.command
    secret_in_request = x_mygcp_secret

    if secret and command.secret != secret and secret != secret_in_request:
        response.status_code = 503
        return Output(
                message=f"Invalid auth. Set {secret_header_key} with valid secret value.",
                return_code=1,
            )

    proc = subprocess.run(cmd_str, shell=True, stdout=PIPE, stderr=PIPE, text=True)
    stdout_output = proc.stdout.split("\n")
    stderr_output = proc.stderr.split("\n")

    return Output(
            stderr_output=stderr_output,
            stdout_output=stdout_output,
            return_code=proc.returncode,
        )

if __name__ == "__main__":
    options = {
            'port': int(port),
            'host': '0.0.0.0',
            'workers': 2,
            'reload': True,
        }
    uvicorn.run("main:app", **options)
