#!/bin/bash

# Define variables
RANDOM_VAR=$(uuidgen| cut -c 1-8)
BUCKET_NAME=skills-korea-2023-${RANDOM_VAR}

# Create new bucket
aws s3 mb --region ap-northeast-2 s3://${BUCKET_NAME}

# Create files
echo "Success!" > public.txt
echo "Failure!" > private.txt

# Upload files
aws s3 --region ap-northeast-2 cp public.txt s3://${BUCKET_NAME}
aws s3 --region ap-northeast-2 cp private.txt s3://${BUCKET_NAME}

# Policy
cat <<EOF > s3_policy.json
{
	"Version": "2012-10-17",
	"Statement": [
		{
			"Sid": "PublicAccess",
			"Effect": "Deny",
			"Principal": "*",
			"Action": "s3:PutObject",
			"Resource": "arn:aws:s3:::${BUCKET_NAME}/public.txt"
		},
		{
			"Sid": "PrivateAccess",
			"Effect": "Allow",
            "Principal": "*",
			"Action": "s3:PutObject",
			"Resource": "arn:aws:s3:::${BUCKET_NAME}/private.txt"
		}
	]
}
EOF
aws s3api put-public-access-block \
    --bucket ${BUCKET_NAME} \
    --public-access-block-configuration "BlockPublicAcls=false,IgnorePublicAcls=false,BlockPublicPolicy=false,RestrictPublicBuckets=false"
aws s3api put-bucket-policy --bucket ${BUCKET_NAME} --policy file://s3_policy.json

# print the created bucket
echo "Bucket Name : $BUCKET_NAME"