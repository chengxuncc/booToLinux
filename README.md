# BooToLinux
Using `bcdedit` to select the next time boot of UEFI: ```bcdedit /set {fwbootmgr} bootsequence {GUID}```

## Download
If you don't want to build by yourself, you can download prebuild binary here: [release](https://github.com/chengxuncc/booToLinux/releases)

## Build 
```dos
go get github.com/bhendo/go-powershell github.com/bhendo/go-powershell/backend
go build -o booToLinux.exe
```
You should right click and run as Administrator.

## Build with Administrator privileges
```dos
go get github.com/bhendo/go-powershell github.com/bhendo/go-powershell/backend
go get github.com/akavel/rsrc
rsrc -manifest booToLinux.exe.manifest -o booToLinux.syso
go build -o booToLinux.exe
```
Double click to run.

## License
BooToLinux is licensed under the MIT license.
