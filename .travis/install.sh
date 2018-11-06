#!/bin/bash

set -x

if [[ $TRAVIS_OS_NAME == 'windows' ]]; then
  # create virtual com using com0com in windows
  export TEST_INPUT_PTY=COM10
  export TEST_OUTPUT_PTY=COM11
  # download com0com with git-bash 
  curl -fsSL https://sourceforge.net/projects/com0com/files/latest/download -o com0com.zip
  7z e com0com.zip *x64*.exe
  mv *x64*.exe com0com-setup.exe
  export CNC_INSTALL_CNCA0_CNCB0_PORTS=YES
  export CNC_INSTALL_COMX_COMX_PORTS=YES
  export CNC_INSTALL_SKIP_SETUP_PREINSTALL=NO
  runas.exe /savecred /user:administrator 'cmd /C com0com-setup.exe /S /D="%ProgramFiles%\com0com"'
  ls -alh "/c/Program Files/"
  ls -alh "/c/Program Files/com0com/"
  ls -alh "/c/Program Files (x86)/"
  ls -alh "/c/Program Files (x86)/com0com/"
  cmd '/C "%ProgramFiles%\com0com\setupc" install Portname=COM10 Portname=COM11'
else
  # create virtual pty
  make pty_start
  sleep 10
  export TEST_INPUT_PTY=$(cat socat.out | cut -d\  -f7 | sed -n 1p)
  export TEST_OUTPUT_PTY=$(cat socat.out | cut -d\  -f7 | sed -n 2p)
fi

set +x
