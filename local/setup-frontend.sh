#!/bin/bash
set -e

cd "$(dirname "$0")"

source "setup-environment.sh"

echo "[INFO] Stopping any existing frontend processes..."

echo "[INFO] Stopping any existing CMS UI processes..."
if [ -d "../frontend/HobbyStreamerCMS" ]; then
  cd ../frontend/HobbyStreamerCMS
  
  pkill -f "npm run web" || true
  pkill -f "expo start" || true
  
  lsof -ti:8081 | xargs kill -9 2>/dev/null || true
  
  cd ../../local
fi

echo "[INFO] Stopping any existing Streaming UI processes..."
if [ -d "../frontend/HobbyStreamerUI" ]; then
  cd ../frontend/HobbyStreamerUI
  
  pkill -f "npm run streaming-ui" || true
  pkill -f "expo start" || true
  
  lsof -ti:8085 | xargs kill -9 2>/dev/null || true
  
  cd ../../local
fi

echo "[INFO] Stopping any existing Expo processes..."
pkill -f "expo start" || true
sleep 2

echo "[INFO] Setting up CMS UI..."
if [ -d "../frontend/HobbyStreamerCMS" ]; then
  cd ../frontend/HobbyStreamerCMS
  
  echo "[INFO] Ensuring correct Node.js version..."
  if command -v nvm &> /dev/null; then
    nvm use
  fi
  
  echo "[INFO] Cleaning and reinstalling CMS UI dependencies..."
  rm -rf node_modules package-lock.json
  npm install
  
  echo "[INFO] Starting CMS UI web application..."
  nohup npm run web > web.log 2>&1 &
  
  echo "[INFO] Waiting for CMS UI web server to start..."
  sleep 10
  
  echo "[INFO] CMS UI web application started"
  echo "[INFO] - CMS UI Web: http://localhost:8081"
  echo "[INFO] To run on device/simulator:"
  echo "[INFO]   - Android: npm run android"
  echo "[INFO]   - iOS: npm run ios"
  
  cd ../../local
else
  echo "[WARNING] CMS UI directory not found at ../frontend/HobbyStreamerCMS"
fi

echo "[INFO] Setting up Streaming UI..."
if [ -d "../frontend/HobbyStreamerUI" ]; then
  cd ../frontend/HobbyStreamerUI

  echo "[INFO] Ensuring correct Node.js version..."
  if command -v nvm &> /dev/null; then
    nvm use
  fi

  echo "[INFO] Cleaning and reinstalling Streaming UI dependencies..."
  rm -rf node_modules package-lock.json
  npm install

  echo "[INFO] Starting Streaming UI web application..."
  nohup npm run streaming-ui > streaming-ui.log 2>&1 &

  echo "[INFO] Waiting for Streaming UI web server to start..."
  sleep 10

  echo "[INFO] Streaming UI web application started"
  echo "[INFO] - Streaming UI Web: http://localhost:8085"
  echo "[INFO] To run on device/simulator:"
  echo "[INFO]   - Android: npm run android"
  echo "[INFO]   - iOS: npm run ios"

  cd ../../local
else
  echo "[WARNING] Streaming UI directory not found at ../frontend/HobbyStreamerUI"
fi

echo "[INFO] Frontend setup completed" 