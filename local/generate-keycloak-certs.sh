#!/bin/bash

# Create certificates directory
mkdir -p keycloak-certs

# Generate self-signed certificate and key
openssl req -x509 -newkey rsa:4096 -keyout keycloak-certs/key.pem -out keycloak-certs/cert.pem -days 365 -nodes -subj "/C=US/ST=State/L=City/O=Organization/CN=localhost"

echo "Self-signed certificates generated in keycloak-certs/"
echo "Certificate: keycloak-certs/cert.pem"
echo "Private Key: keycloak-certs/key.pem" 