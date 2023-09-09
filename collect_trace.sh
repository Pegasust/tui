#!/usr/bin/env sh
# Bash script to aggregate trace data
for file in trace/*.log; do
  # Analyze or merge logic here
  echo "# $file"
  go tool trace $file
done
