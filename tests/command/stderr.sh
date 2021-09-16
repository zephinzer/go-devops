#!/bin/bash
#      ^ bash is required because of read -t ${SECONDS} x;

>&2 echo "this prints to stderr";
read -t 1 x;
>&2 echo "you entered: $x";
