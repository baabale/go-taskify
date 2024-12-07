#!/bin/bash

# Exit on any error
set -e

echo "Starting deployment process..."

# Update system packages
echo "Updating system packages..."
sudo yum update -y

# Install Docker
echo "Installing Docker..."
sudo yum install -y docker
sudo service docker start
sudo usermod -a -G docker ec2-user

# Install Docker Compose
echo "Installing Docker Compose..."
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Create app directory
echo "Setting up application directory..."
mkdir -p ~/taskify
cd ~/taskify

# Copy application files
echo "Copying application files..."
# Note: These files should be copied from your local machine to the EC2 instance
# using scp or similar before running this script

# Set correct permissions
echo "Setting permissions..."
sudo chown -R ec2-user:ec2-user ~/taskify
chmod -R 755 ~/taskify

# Start the containers
echo "Starting Docker containers..."
sudo docker-compose down --volumes --remove-orphans || true
sudo docker-compose up --build -d

# Show container status
echo "Checking container status..."
sudo docker-compose ps

# Show application logs
echo "Showing application logs..."
sudo docker-compose logs app

echo "Deployment complete! The application should be running at http://YOUR_EC2_IP:3000"
echo "Don't forget to configure your EC2 security group to allow traffic on port 3000"
