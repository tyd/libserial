#!/bin/bash

set -x

if [[ "$TRAVIS_OS_NAME" != "windows" ]]; then 
  if [[ "$(stty speed < ${TEST_OUTPUT_PTY})" != "${TEST_BAUD_RATE}" ]] then
    exit 1
  fi
fi
