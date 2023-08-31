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
