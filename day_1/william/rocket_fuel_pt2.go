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
        tFuel += recursiveFuel(m) //go floors normally
    }
    fmt.Println(tFuel)
}

func recursiveFuel(mass int) int {
    f := 0
    for{
        mass = fuelCalc(mass)
        f += mass
        if mass <= 0 {
            return f
        }
    }
}

func fuelCalc(val int) int{
    x := val/3 -2
    if x < 0 {
        x = 0
    }
    return x
}




func errCheck(e error){
    if e!= nil{
        log.Fatalln(e)
    }
}
