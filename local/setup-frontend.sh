#!/bin/bash
set -e

cd "$(dirname "$0")"

source ./setup-environment.sh

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

echo "[INFO] Frontend setup completed" 