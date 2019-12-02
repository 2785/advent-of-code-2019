package main

import (
    "bufio"
    "fmt"
    "log"
    "os"
    "strconv"
)

func main(){
    file,err:= os.Open("./in")
    errCheck(err)
    defer file.Close()
    s := bufio.NewScanner(file)
    tFuel := 0
    for s.Scan(){
        m,err := strconv.Atoi(s.Text())
        errCheck(err)
        tFuel += m/3 -2 //go floors normally
    }
    fmt.Println(tFuel)
}

func errCheck(e error){
    if e!= nil{
        log.Fatalln(e)
    }
}
