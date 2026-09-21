#!/usr/bin/env bash
set -euo pipefail

base_url=${1:?usage: openai-chat-smoke.sh <base-url> <model>}
model=${2:?usage: openai-chat-smoke.sh <base-url> <model>}

chat_payload=$(jq -nc --arg model "$model" '{model:$model,messages:[{role:"user",content:"reply with ok"}],max_tokens:4}')
completion_payload=$(jq -nc --arg model "$model" '{model:$model,prompt:"reply with ok",max_tokens:4}')

chat_response=$(mktemp)
completion_response=$(mktemp)
cleanup() {
	rm -f "$chat_response" "$completion_response"
}
trap cleanup EXIT

if curl -fsS "${base_url%/}/v1/chat/completions" \
	-H 'content-type: application/json' \
	-d "$chat_payload" >"$chat_response"; then
	jq -e '.choices[0].message.content or .choices[0].text' "$chat_response"
	exit 0
fi

curl -fsS "${base_url%/}/v1/completions" \
	-H 'content-type: application/json' \
	-d "$completion_payload" >"$completion_response"
jq -e '.choices[0].text or .choices[0].message.content' "$completion_response"
