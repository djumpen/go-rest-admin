#!/bin/bash

target_folder="postman"
postman_api_key="123456"
postman_collection_id="78901"

mkdir -p $target_folder

curl -X GET https://api.getpostman.com/collections/$postman_collection_id -H "X-Api-Key: $postman_api_key" -H "Cache-Control: no-cache" -o $target_folder/postman_collection.json