#!/bin/bash

set -x

if [[ $TRAVIS_OS_NAME == 'windows' ]]; then
  # create virtual com using com0com in windows

  # download com0com with git-bash 
  curl -fsSL https://sourceforge.net/projects/com0com/files/latest/download -o com0com.zip
  7z e com0com.zip *x64*.exe
  mv *x64*.exe com0com-setup.exe
  # use cmd.exe to run install.bat
  powershell -c start -verb runas cmd '/C .travis\install.bat /D "%V"'
  export TEST_INPUT_PTY=COM10
  export TEST_OUTPUT_PTY=COM11
else
  # create virtual pty
  make pty_start
  sleep 10
  export TEST_INPUT_PTY=$(cat socat.out | cut -d\  -f7 | sed -n 1p)
  export TEST_OUTPUT_PTY=$(cat socat.out | cut -d\  -f7 | sed -n 2p)
fi

set +x
