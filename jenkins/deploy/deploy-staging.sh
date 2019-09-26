#!/usr/bin/env bash
# Staging environment: 31.171.247.162
# Private key for ssh: /opt/keypairs/ditas-testbed-keypair.pem

# TODO state management? We are killing without caring about any operation the conainer could be doing.

ssh -i /opt/keypairs/ditas-testbed-keypair.pem cloudsigma@31.171.247.162 << 'ENDSSH'
# Ensure that a previously running instance is stopped (-f stops and removes in a single step)
# || true - "docker stop" failt with exit status 1 if image doesn't exists, what makes the Pipeline fail. the "|| true" forces the command to exit with 0.

sudo docker stop --time 20 computation-movement-enactor || true
sudo docker rm --force computation-movement-enactor || true
sudo docker pull ditas/computation-movement-enactor:latest

sudo docker run -p 30090:8080 -d --name computation-movement-enactor ditas/computation-movement-enactor:latest
ENDSSH