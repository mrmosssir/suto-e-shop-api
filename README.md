# suto-e-shop-api

## 部署到 gcloud run
gcloud run deploy "TARGET" --source . --region "REGION" --project "PROJECT_NAME" --allow-unauthenticated --clear-base-image

## 覆蓋 config 設定檔
gcloud run services replace config.yaml --region="REGION" --project="PROJECT_NAME"

## 設定 storage bucket 權限可被用戶訪問
gcloud storage buckets add-iam-policy-binding "YOUR_FIREBASE_STORAGE" \
  --member="allUsers" \
  --role="roles/storage.objectViewer"
