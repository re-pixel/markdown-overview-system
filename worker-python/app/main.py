import boto3
import requests
import json
import time

sqs = boto3.client("sqs", endpoint_url="http://localhost:4566", region_name="eu-central-1")
s3 = boto3.client("s3", endpoint_url="http://localhost:4566", region_name="eu-central-1")

TASK_QUEUE_URL ="http://sqs.eu-central-1.localhost.localstack.cloud:4566/000000000000/task-queue"
RESPONSE_QUEUE_URL = "http://sqs.eu-central-1.localhost.localstack.cloud:4566/000000000000/response-queue"
OPENROUTER_API_KEY = "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"
OPENROUTER_URL = "https://openrouter.ai/api/v1/chat/completions"

def process_file(bucket, key):

    obj = s3.get_object(Bucket=bucket, Key=key)
    file_content = obj["Body"].read().decode("utf-8")

    file_content = "Summarize following text in two sentences:\n" + file_content

    headers = {
        "Authorization": f"Bearer {OPENROUTER_API_KEY}",
        "Content-Type": "application/json",
    }

    data = {
        "model": "x-ai/grok-4-fast:free",  
        "messages": [{"role": "user", "content": file_content}]
    }

    response = requests.post(OPENROUTER_URL, headers=headers, json=data)
    resp_json = response.json()
    print("OpenRouter response:", resp_json)

    overview = resp_json["choices"][0]["message"]["content"]

    extension_pos = key.index(".")
    overview_key = key[:extension_pos] + "_overview.txt"

    s3.put_object(
        Bucket=bucket,
        Key=overview_key,
        Body=overview.encode("utf-8")  
    )
    print(f"Summary uploaded to s3://{bucket}/{overview_key}")
    return overview_key


def worker_loop():
    while True:
        resp = sqs.receive_message(
            QueueUrl=TASK_QUEUE_URL,
            MaxNumberOfMessages=5,
            WaitTimeSeconds=10
        )

        messages = resp.get("Messages", [])
        if not messages:
            print("No new messages, waiting...")
            time.sleep(5)
            continue

        for msg in messages:
            body = json.loads(msg["Body"])
            bucket = body["bucket"]
            key = body["key"]
            userId = body["userId"]

            print(f"Processing file {key} from {bucket} for user {userId}")
            overview_key = process_file(bucket, key)
            response_msg = {
                "bucket": bucket,
                "key": overview_key,
                "status": "completed",
                "userId": userId,
            }

            sqs.send_message(
                QueueUrl=RESPONSE_QUEUE_URL,
                MessageBody=json.dumps(response_msg),
            )
            print(f"Sent response message for {overview_key}")


            sqs.delete_message(
                QueueUrl=TASK_QUEUE_URL,
                ReceiptHandle=msg["ReceiptHandle"]
            )
            print(f"Deleted message {msg['MessageId']}")


if __name__ == "__main__":
    print("Worker started...")
    worker_loop()