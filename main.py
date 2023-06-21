import os
from fastapi import FastAPI, Depends, Header, Request, APIRouter
from pydantic import *
import uvicorn
import json
import subprocess
from subprocess import PIPE

app = FastAPI()

secret = os.environ.get("SECRET")
port = os.environ.get("PORT", "8080")
SECRET_HEADER_KEY = "X-MyGCP-Secret"

class Command(BaseModel):
    command: str = ""
    secret: str = ""

@app.post("/shellcommand")
def cmd(command: Command, request: Request, x_mygcp_secret = Header(default=None)):

    cmd_str = command.command
    secret_in_request = x_mygcp_secret

    if secret and command.secret != secret and secret != secret_in_request:
        return {"output":f"Invalid auth. Set {SECRET_HEADER_KEY} with valid secret value."}

    proc = subprocess.run(cmd_str, shell=True, stdout=PIPE, stderr=PIPE, text=True)
    stdout_output = proc.stdout.split("\n")
    stderr_output = proc.stderr.split("\n")

    return json.dumps({"stdout_output":stdout_output, "stderr_output":stderr_output, "return_code":proc.returncode})

if __name__ == "__main__":
    # app.run(debug=True, host="0.0.0.0", port=int(os.environ.get("PORT", 8080)))
    options = {
            'port': int(port),
            'host': '0.0.0.0',
            'workers': 2,
            'reload': True,
        }
    uvicorn.run("main:app", **options)
