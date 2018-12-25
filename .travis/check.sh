#!/bin/bash

set -x

if [[ $TRAVIS_OS_NAME == 'windows' ]]; then
  powershell ./.travis/port_status.ps1
else
  if [[ "$(stty speed < ${TEST_OUTPUT_PTY})" != "${TEST_BAUD_RATE}" ]]; then
    exit 1
  fi
fi
