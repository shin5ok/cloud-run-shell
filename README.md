## Warning
### Don't use --allow-authenticated option while deploying your Cloud Run
### You have to use this api carefully or it would be serious risky

## Usage
```
V=`gcloud auth print-identity-token`
URL={URL of Cloud Run}
echo '{"secret":"your-secret", "command":"curl http://10.146.0.2"}' | curl -d @- -H "Authorization: Bearer $V" $URL/shellcommand
```