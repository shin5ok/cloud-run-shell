SECRET=${SECRET:-gcp}
GEN=${GEN:-gen2}
SERVICE_NAME=${SERVICE_NAME:-shellcommand}

gcloud run deploy ${SERVICE_NAME} --source=. \
    --set-env-vars=SECRET=${SECRET} \
    --region=asia-northeast1 \
    --execution-environment=${GEN} \
    --memory=1Gi \
    --cpu=1 \
    --max-instances=1 \
    --timeout=3600 \
    $@
