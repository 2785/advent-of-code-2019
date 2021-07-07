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
    arOld := make([]int,len(ar))
    copy(arOld,ar)
    for n:=0; n<100; n++{
        for v :=0; v<100; v++{
            x := execute(ar,n,v)
            if x == 19690720{
                fmt.Println(100*n+v)
                return
            }
            copy(ar, arOld)
        }
    }
}

func execute(ar []int, n int, v int) int{
    ar[1] = n
    ar[2] = v
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
            return ar[0]
        }
    }

}

func errCheck(e error){
    if e!= nil{
        log.Fatalln(e)
    }
}

