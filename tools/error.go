package tools

import "log"

func PrintError(err error){
	log.Printf("Error: %+v", err)
}

func Error(err error,s string){
	log.Print(err)
	log.Println(s)
}

//func奥术大师git