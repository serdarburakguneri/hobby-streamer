#!/bin/bash

set -e

echo "Running all backend tests..."
echo "================================"

# Function to run tests in a directory
run_tests() {
    local dir=$1
    local service_name=$2
    
    if [ -d "$dir" ]; then
        echo "Testing $service_name..."
        cd "$dir"
        
        if [ -f "go.mod" ]; then
            go test ./... -v
            if [ $? -eq 0 ]; then
                echo "$service_name tests passed"
            else
                echo "$service_name tests failed"
                exit 1
            fi
        else
            echo "No go.mod found in $service_name, skipping..."
        fi
        
        cd - > /dev/null
        echo ""
    else
        echo "Directory $dir not found, skipping $service_name..."
        echo ""
    fi
}

# Run tests for each service
run_tests "asset-manager" "Asset Manager"
run_tests "auth-service" "Auth Service"
run_tests "streaming-api" "Streaming API"
run_tests "transcoder" "Transcoder"
run_tests "pkg" "Shared Packages"

# Run tests for lambdas
echo "Testing Lambdas..."
cd lambdas
for lambda_dir in */; do
    if [ -d "$lambda_dir" ]; then
        lambda_name=$(basename "$lambda_dir")
        echo "Testing Lambda: $lambda_name"
        cd "$lambda_dir"
        
        if [ -f "go.mod" ]; then
            go test ./... -v
            if [ $? -eq 0 ]; then
                echo "Lambda $lambda_name tests passed"
            else
                echo "Lambda $lambda_name tests failed"
                exit 1
            fi
        else
            echo "No go.mod found in lambda $lambda_name, skipping..."
        fi
        
        cd ..
        echo ""
    fi
done
cd ..

echo "================================"
echo "All tests completed successfully!" 