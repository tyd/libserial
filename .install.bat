set CNC_INSTALL_CNCA0_CNCB0_PORTS=YES
set CNC_INSTALL_COMX_COMX_PORTS=YES
set CNC_INSTALL_SKIP_SETUP_PREINSTALL=NO
com0com-setup.exe /S
cd "C:\Program Files (x86)\com0com"
setupc install Portname=COM10 Portname=COM11
