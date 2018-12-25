$COMportList = [System.IO.Ports.SerialPort]::getportnames()

ForEach ($COMport in $COMportList) {
    $port = new-object System.IO.Ports.SerialPort $COMport
    echo $port.BaudRate
    $port.Dispose()
    echo ";"
}
