#!/bin/bash

set -x

if [[ $TRAVIS_OS_NAME == 'windows' ]]; then
  # create virtual com using com0com in windows
  export TEST_INPUT_PTY=COM10
  export TEST_OUTPUT_PTY=COM11
  # install options
  export CNC_INSTALL_CNCA0_CNCB0_PORTS=YES
  export CNC_INSTALL_COMX_COMX_PORTS=YES
  export CNC_INSTALL_SKIP_SETUP_PREINSTALL=NO
  # install com0com via choco
  (choco install --force --yes -v com0com) || cat /c/ProgramData/chocolatey/logs/chocolatey.log
  ls -alh ${ProgramFiles}
  ls -alh ${ProgramFiles(x86)}
  cmd '/C "%ProgramFiles(x86)%\com0com\setupc" install Portname=COM10 Portname=COM11'
else
  # create virtual pty
  make pty_start
  sleep 10
  export TEST_INPUT_PTY=$(cat socat.out | cut -d\  -f7 | sed -n 1p)
  export TEST_OUTPUT_PTY=$(cat socat.out | cut -d\  -f7 | sed -n 2p)
fi

set +x
