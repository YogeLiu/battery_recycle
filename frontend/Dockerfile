# ---- Build Stage ----
FROM node:20-alpine AS builder

WORKDIR /app

# Set npm registry for faster downloads
RUN npm config set registry https://registry.npmmirror.com

# Copy package files
COPY package*.json ./
RUN npm install

# Copy source code and build
COPY . .
RUN npm run build

# ---- Production Stage ----
FROM nginx:stable-alpine

# Copy build files to nginx
COPY --from=builder /app/build /usr/share/nginx/html

# Expose port
EXPOSE 80

# Start nginx
CMD ["nginx", "-g", "daemon off;"]