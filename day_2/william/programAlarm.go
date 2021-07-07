package main

import (
    "fmt"
    "log"
    "io/ioutil"
    "strconv"
    "strings"
)

func main(){
    file,err:=ioutil.ReadFile("./in")
    errCheck(err)
    str := string(file)
    str = strings.TrimSuffix(str, "\n")
    s := strings.Split(str,",")
    ar :=make([]int, len(s))
    for i,v := range s{
        ar[i], err = strconv.Atoi(v)
        errCheck(err)
    }
    ar[1] = 12
    ar[2] = 2

    i := 0
    for{
        if ar[i] == 1{
            ar[ar[i+3]] = ar[ar[i+1]] + ar[ar[i+2]]
            i += 4
        }
        if ar[i] == 2{
            ar[ar[i+3]] = ar[ar[i+1]] * ar[ar[i+2]]
            i += 4
        }
        if ar[i] == 99{
            fmt.Println(ar[0])
            return
        }
    }

}

func errCheck(e error){
    if e!= nil{
        log.Fatalln(e)
    }
}

