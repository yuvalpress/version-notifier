#!/bin/bash

#Define the basic variables
if [ -z "$PROFILE" ]; then
    PROFILE="terraform"
fi
if [ -z "$REGION" ]; then
    REGION="eu-west-1"
fi
# Set your AWS S3 bucket and key
S3_BUCKET="sirrend-cloud-prod-lambda"
ZIP_NAME="version_notifier.zip"
LAMBDA_NAME="sirrend-cloud-prod-version-notifier"

# Create a temporary directory
mkdir -p tmp
TMP_DIR="./tmp"

# Copy all Go files to the temporary directory
cp *.go $TMP_DIR/

# Download Go dependencies (optional, remove if not needed)
go get ./...

# Build the Go binary inside the temporary directory
GOOS=linux GOARCH=arm64 go build -tags lambda.norpc -o $TMP_DIR/bootstrap $TMP_DIR/main.go

# Change to the temporary directory
cd $TMP_DIR

# Create a ZIP archive containing all Go files and the binary
zip -r $ZIP_NAME .

# Upload the ZIP archive to S3
aws s3 cp $ZIP_NAME s3://$S3_BUCKET/

# Update the function lambda with the latest artifact
aws lambda update-function-code --function-name $LAMBDA_NAME --s3-bucket $S3_BUCKET --s3-key $ZIP_NAME --profile $PROFILE --region $REGION > /dev/null

# Clean up temporary files
cd -
rm -rf $TMP_DIR
