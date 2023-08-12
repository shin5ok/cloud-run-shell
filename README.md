## Warning
### Don't use --allow-authenticated option while deploying your Cloud Run
### You have to use this api carefully or it would be serious risky

## Usage
### Prepare to use in common
```
export SECRET=<your secret>
export URL=<url of Cloud Run service>
```
If the Cloud Run sevice requires ID Token,
```
export TOKEN=$(gcloud auth print-identity-token)
```

### With CLI
1. Build it to make client as 'client/shell' once.
```
make client
```
2. Request
```
cd client/
./shell ps aux | jq .
```

### As REST with curl
```
echo '{"command":"ps aux"}' | curl -s -d @- -H "X-MyGCP-Secret: $SECRET" -H "Authorization: Bearer $TOKEN" $URL/shellcommand | jq .
```
