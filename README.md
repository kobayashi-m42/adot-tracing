# adot-tracing
## Description
AWS Distro for OpenTelemetryを利用してトレースをAWS X-Rayで可視化するサンプルアプリケーションです。
以下のエンドポイントにリクエストするとトレースがAWS X-Rayに送信されます。

- /xray-sdk/outgoing-http-call
    - X-Ray SDKを利用してaws.amazon.com（http://aws.amazon.com/ ）にHTTPリクエストを行う
- /xray-sdk/aws-sdk-call
    - X-Ray SDKを利用してS3バケットをリストアップする
- /otel-sdk/outgoing-http-call
    - OpenTelemetry SDKを利用して aws.amazon.com（http://aws.amazon.com/ ）にHTTPリクエストを行う
- /otel-sdk/aws-sdk-call
    - OpenTelemetry SDKを利用してS3バケットをリストアップする

## Requirements
Go 1.20+

Docker 20.10.22

## Getting Started
AWS のクレデンシャル情報を環境変数`AWS_ACCESS_KEY_ID`,`AWS_SECRET_ACCESS_KEY`に設定してください。

以下のコマンドを実行し、Docker上でアプリケーションを実行します。

```shell
docker compose up --build -d
```

以下のエンドポイントにリクエストすることでAWS X-Rayでトレースが送信されます。
```shell
curl http://localhost:8000/xray-sdk/outgoing-http-call
curl http://localhost:8000/xray-sdk/aws-sdk-call
curl http://localhost:8000/otel-sdk/outgoing-http-call
curl http://localhost:8000/otel-sdk/aws-sdk-call
```
