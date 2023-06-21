# Copyright 2020 Google, LLC.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# [START cloudrun_helloworld_service]
# [START run_helloworld_service]
import os

from flask import Flask, render_template, request

app = Flask(__name__)

secret = os.environ.get("SECRET")
SECRET_HEADER_KEY = "X-MyGCP-Secret"


@app.route("/")
def hello_world():
    name = os.environ.get("NAME", "World")
    return "Hello {}!".format(name)

@app.route("/shellcommand", methods=['POST'])
def cmd():
    import subprocess
    from subprocess import PIPE
    data = request.get_data()
    secret_in_request = request.headers.get(SECRET_HEADER_KEY, "")
    import json, os
    try:
        j = json.loads(data)
        print(j)
    except Exception as e:
        print(str(e))
        return {"output":f"Error {str(e)}"}
    if secret and j.get('secret') != secret and secret != secret_in_request:
        return {"output":f"Invalid auth. Set {SECRET_HEADER_KEY} with valid secret value."}
    proc = subprocess.run(j['command'], shell=True, stdout=PIPE, stderr=PIPE, text=True)
    stdout_output = proc.stdout.split("\n")
    stderr_output = proc.stderr.split("\n")
    print(stdout_output)

    return json.dumps({"stdout_output":stdout_output, "stderr_output":stderr_output, "return_code":proc.returncode})

if __name__ == "__main__":
    app.run(debug=True, host="0.0.0.0", port=int(os.environ.get("PORT", 8080)))
# [END run_helloworld_service]
# [END cloudrun_helloworld_service]
