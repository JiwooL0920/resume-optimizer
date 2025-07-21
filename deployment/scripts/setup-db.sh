#!/bin/bash

# Setup database for resume optimizer
# This script creates the database and runs migrations on your CNPG cluster

set -e

NAMESPACE="cnpg-system"
CLUSTER_NAME="postgresql-cluster"
DB_NAME="resume_optimizer"

echo "Creating database $DB_NAME in CNPG cluster..."

# Create database
kubectl exec -n $NAMESPACE $CLUSTER_NAME-1 -- psql -U postgres -c "CREATE DATABASE $DB_NAME;"

echo "Running migrations..."

# Run migrations directly
kubectl exec -n $NAMESPACE $CLUSTER_NAME-1 -- psql -U postgres -d $DB_NAME -c "$(cat shared/database/migrations/001_initial.sql)"

echo "Database setup complete!"
echo "Connection string: postgres://postgres:password@postgresql-cluster-rw.cnpg-system.svc.cluster.local:5432/$DB_NAME"