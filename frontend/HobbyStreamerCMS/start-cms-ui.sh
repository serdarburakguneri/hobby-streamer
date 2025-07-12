#!/bin/bash

# Ensure correct Node.js version
if command -v nvm &> /dev/null; then
  echo "Using nvm to set Node.js version..."
  nvm use
fi

# Check Node.js version
NODE_VERSION=$(node --version)
echo "Using Node.js version: $NODE_VERSION"

# Start the web application
echo "Starting Expo web application..."
npm run web 