#!/bin/bash
if [ "$1" = "get" ]; then
  TOKEN=$(curl -s http://localhost:3100/get-token)
  echo "username=will"
  echo "password=$TOKEN"
fi
