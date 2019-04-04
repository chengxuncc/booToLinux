package main

import (
	"fmt"
	"github.com/bhendo/go-powershell"
	"github.com/bhendo/go-powershell/backend"
	"log"
	"os"
	"strings"
)

func pauseFatal(output string){
	fmt.Println(output)
	fmt.Print(`Enter to exit:`)
	fmt.Scanln()
	os.Exit(0)
}

func main(){
	shell, err := powershell.New(&backend.Local{})
	if err != nil {
		log.Panicln("Open Powershell failed:", err.Error())
	}
	defer shell.Exit()
	stdout,_,err:=shell.Execute("bcdedit /enum firmware")
	if err!=nil{
		log.Panicln("Execute Error:", err.Error())
	}
	windowsBootManager:="{bootmgr}"
	if !strings.Contains(stdout,windowsBootManager){
		pauseFatal(stdout)
	}

	firmwareLines:=strings.Split(stdout,"\r\n")
	UEFIBootList:=make(map[string]string)
	var identifierList []string
	currentIdentifier :=""
	for _,line:=range firmwareLines{
		if len(line)==0 || line[0]=='-'{
			continue
		}
		lineSplit:=strings.Split(line,`  `)
		switch lineSplit[0] {
		case "identifier":
			currentIdentifier =lineSplit[len(lineSplit)-1]
			identifierList=append(identifierList,currentIdentifier)
		case "description":
			description:=lineSplit[len(lineSplit)-1]
			description=strings.Trim(description,` `)
			UEFIBootList[currentIdentifier]=description
		}
	}

	fmt.Println("----UEFI Boot List----")
	for k,v:=range identifierList{
		if k==0{
			continue
		}
		fmt.Println(k,UEFIBootList[v])
	}

	var inputIndex int
	defaultIndex:=1
	for k,v:=range identifierList{
		if k==0{
			continue
		}
		if strings.Contains(UEFIBootList[v],"Windows Boot Manager"){
			continue
		}
		if strings.Contains(UEFIBootList[v],"Internal"){
			continue
		}
		defaultIndex=k
		break
	}
	INPUT:
	fmt.Printf("\nPlease input index of the next boot, Enter [%s]:",UEFIBootList[identifierList[defaultIndex]])
	_,err=fmt.Scanln(&inputIndex)
	if err!=nil{
		switch err.Error() {
		case "unexpected newline":
			inputIndex=defaultIndex
			break
		case "EOF":
			os.Exit(0)
		default:
			fmt.Println(err.Error())
			goto INPUT
		}
	}
	if inputIndex>=len(identifierList){
		fmt.Println("Invalid Index")
		goto INPUT
	}

	setNextBoot(shell,identifierList[inputIndex])
	fmt.Println("Next boot:",UEFIBootList[identifierList[inputIndex]])

	fmt.Print("Enter to reboot:")
	_,err=fmt.Scanln(&inputIndex)
	if err!=nil && err.Error()=="unexpected newline"{
		rebootNow(shell)
	}

}

func setNextBoot(shell powershell.Shell,identifier string){
	cmd:=fmt.Sprintf(`bcdedit /set "{fwbootmgr}" bootsequence "%s"`,identifier)
	stdout,_,err:=shell.Execute(cmd)
	if err!=nil{
		pauseFatal(err.Error())
	}
	fmt.Print(stdout)
}

func rebootNow(shell powershell.Shell){
	cmd:=`shutdown /r /t 0`
	_,_,err:=shell.Execute(cmd)
	if err!=nil{
		pauseFatal(err.Error())
	}
}
