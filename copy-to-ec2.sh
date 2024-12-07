#!/bin/bash

# Check if EC2 instance IP is provided
if [ -z "$1" ]; then
    echo "Usage: $0 <ec2-ip> [pem-file]"
    echo "Example: $0 12.34.56.78 EC2TASKIFY.pem"
    exit 1
fi

EC2_IP=$1
PEM_FILE="${2:-EC2TASKIFY.pem}"  # Default to EC2TASKIFY.pem if not provided

# Check if PEM file exists
if [ ! -f "$PEM_FILE" ]; then
    echo "PEM file not found: $PEM_FILE"
    exit 1
fi

echo "Creating temporary directory for deployment..."
# Create deployment directory if it doesn't exist
ssh -i "$PEM_FILE" -o StrictHostKeyChecking=no "ec2-user@$EC2_IP" "mkdir -p ~/taskify"

# Copy required files to EC2
echo "Copying files to EC2 instance..."
scp -i "$PEM_FILE" \
    Dockerfile \
    docker-compose.yml \
    go.mod \
    go.sum \
    main.go \
    deploy.sh \
    "ec2-user@$EC2_IP:~/taskify/"

# Copy directories
for dir in config controllers middleware models routes utils; do
    if [ -d "$dir" ]; then
        echo "Copying $dir directory..."
        scp -i "$PEM_FILE" -r "$dir" "ec2-user@$EC2_IP:~/taskify/"
    fi
done

# Make deploy.sh executable and run it
echo "Running deploy script on EC2..."
ssh -i "$PEM_FILE" "ec2-user@$EC2_IP" "cd ~/taskify && chmod +x deploy.sh && ./deploy.sh"

echo "Deployment complete! Your application should be running at http://$EC2_IP:3000"
