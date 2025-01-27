#!/usr/bin/env bash

# Get a list of all running Docker containers
containers=$(docker ps -a -q)

# Check if the list is not empty
if [ -n "$containers" ]; then
  printf "Found Docker containers on the machine.\nAttempting to stop and remove all running Docker containers.\n"
  # Attempt to stop all running containers
  if docker stop $containers; then
      echo "All running Docker containers stopped successfully."
  else
    echo "Error: failed to stop all running Docker containers."
  fi
  # Attempt to remove all containers
  if docker rm $containers; then
    echo "All available Docker containers removed successfully."
  else
    echo "Error: failed to remove available Docker containers."
  fi
else
  echo "No Docker containers found."
fi

docker container prune -f
docker image prune -a -f
docker volume prune -f
