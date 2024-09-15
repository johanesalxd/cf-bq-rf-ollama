BigQuery Remote Function (BQ RF) with Ollama
-----------------------------
Details TBA

# How to run
## Run locally
```
FUNCTION_TARGET=BQRFOllama PROJECT_ID="YOUR_PROJECT_ID" LOCATION="YOUR_LOCATION" CONCURRENCY_LIMIT="YOUR_CONCURRENCY_LIMIT" OLLAMA_URL="YOUR_OLLAMA_URL" go run cmd/main.go
```

## Run locally with Pack and Docker
```
pack build --builder=gcr.io/buildpacks/builder cf-bq-rf-ollama

gcloud auth application-default login

ADC=~/.config/gcloud/application_default_credentials.json && \
docker run -p8080:8080 \
-e PROJECT_ID="YOUR_PROJECT_ID" \
-e LOCATION="YOUR_LOCATION" \
-e CONCURRENCY_LIMIT="YOUR_CONCURRENCY_LIMIT" \
-e OLLAMA_URL="YOUR_OLLAMA_URL" \
-e GOOGLE_APPLICATION_CREDENTIALS=/tmp/keys/secret.json \
-v ${ADC}:/tmp/keys/secret.json \
cf-bq-rf-ollama
```

## Test locally (accept [BQ RF request contract](https://cloud.google.com/bigquery/docs/remote-functions#input_format))
```
curl -m 60 -X POST localhost:8080 \
-H "Content-Type: application/json" \
-d '{
  "requestId": "",
  "caller": "",
  "sessionUser": "",
  "userDefinedContext": {},
  "calls": [
    ["what is bigquery", "gemma2:9b"],
    ["no model found", ""],
    ["error"]
  ]
}'
```

## Run on Cloud Function
```
gcloud functions deploy cf-bq-rf-ollama \
    --gen2 \
    --concurrency=8 \
    --cpu=1 \
    --memory=512Mi \
    --runtime=go122 \
    --region=us-central1 \
    --source=. \
    --entry-point=BQRFOllama \
    --trigger-http \
    --allow-unauthenticated \
    --env-vars-file=.env.yaml
```

## Run on Cloud Run
[![Run on Google Cloud](https://deploy.cloud.run/button.svg)](https://deploy.cloud.run)

# Additional notes
TBA

## Related links
* https://cloud.google.com/bigquery/docs/remote-functions
* https://cloud.google.com/functions/docs/concepts/go-runtime
* https://cloud.google.com/docs/buildpacks/build-function
* https://cloud.google.com/run/docs/tutorials/gpu-gemma2-with-ollama
