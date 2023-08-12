SECRET=${SECRET:-gcp}

gcloud run deploy shellcommand --source=. \
    --set-env-vars=SECRET=${SECRET} \
    --region=asia-northeast1 \
    --execution-environment=gen2 \
    --memory=1Gi \
    --cpu=1 \
    --max-instances=1 \
    --timeout=3600 \
    $@
