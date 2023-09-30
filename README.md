## Warning
### Don't use --allow-authenticated option while deploying your Cloud Run
### You have to use this api carefully or it would be serious risky

## Usage
### Prepare to use in common
```
export SECRET=<your secret>
export URL=<url of Cloud Run service>
```
If the Cloud Run sevice requires ID Token, run as below,
```
export TOKEN=$(gcloud auth print-identity-token)
```

### With CLI

#### 1. Build it to make client as 'client/shell' once.
```
make client
```
Change directory to 'client'.
```
cd client/
```
You're ready to run something on your Cloud Run service.

#### 2. Just Run it

Run it with your any command.
```
./shell ps aux
```

You may see output formatted JSON as below.
```
JSON_MODE=1 ./shell ps aux
```

### Option: Simple pseudo shell without typing client command

After this, You should input command you want in one liner.
```
xargs -L1 ./client/shell
```

Like this,
```
$ xargs -L1 ./client/shell
ls -l
total 44
-rw-r--r-- 1 root root  1109 Sep 30 02:31 Dockerfile
-rw-r--r-- 1 root root    31 Sep 15 14:06 Procfile
drwxr-xr-x 2 root root     0 Sep 30 02:33 byte-unixbench
-rw-r--r-- 1 root root   321 Sep 15 14:06 deploy.sh
drwxr-xr-x 2 root root    80 Sep 30 02:49 google-cloud-sdk
-rw-r--r-- 1 root root  2472 Sep 27 09:51 main.py
-rw-r--r-- 1 root root 39332 Sep 15 14:06 poetry.lock
-rw-r--r-- 1 root root   444 Sep 15 14:06 pyproject.toml

gcloud storage cp gs://shingo-ar-test0729/testvideo.mp4 /tmp

ls /tmp/
cloudsql-proxy-tmp
testvideo.mp4

cd /tmp; ffmpeg -i testvideo.mp4 testvideo.mov

ls -l /tmp/testvideo.mov
-rw-r--r-- 1 root root 29605125 Sep 30 02:54 /tmp/testvideo.mov
```
