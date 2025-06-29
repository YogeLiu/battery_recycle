#!/bin/bash

# Navigate to frontend directory
cd frontend

# Install dependencies if node_modules doesn't exist
if [ ! -d "node_modules" ]; then
    echo "Installing frontend dependencies..."
    npm install
fi

# Build the frontend for production
echo "Building React frontend..."
npm run build

# Serve the built files (you can also use npm run preview)
echo "Frontend built successfully. Files are in frontend/dist/"
echo "You can serve them with: npm run preview"
echo "Or deploy the dist/ folder to your web server."