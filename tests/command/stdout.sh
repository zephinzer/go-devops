#!/bin/bash
#      ^ bash is required because of read -t ${SECONDS} x;

echo "this prints to stdout";
read -t 1 x;
echo "you entered: $x"
