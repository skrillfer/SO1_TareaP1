package main

import (
    "fmt"
    "io/ioutil"
    "strconv"
    "strings"
    "time"
    "log"
)

func getCPUSample() (idle, total uint64) {
    contents, err := ioutil.ReadFile("/proc/stat")
    if err != nil {
        return
    }
    lines := strings.Split(string(contents), "\n")
    for _, line := range(lines) {
        fields := strings.Fields(line)
        if fields[0] == "cpu" {
            numFields := len(fields)
            for i := 1; i < numFields; i++ {
                val, err := strconv.ParseUint(fields[i], 10, 64)
                if err != nil {
                    fmt.Println("Error: ", i, fields[i], err)
                }
                total += val // tally up all the numbers to get total ticks
                if i == 4 {  // idle is the 5th field in the cpu line
                    idle = val
                }
            }
            return
        }
    }
    return
}

func main() {
    
    idle0, total0 := getCPUSample()
    time.Sleep(3 * time.Second)
    idle1, total1 := getCPUSample()

    idleTicks := float64(idle1 - idle0)
    totalTicks := float64(total1 - total0)
    cpuUsage := 100 * (totalTicks - idleTicks) / totalTicks

    fmt.Printf("El uso del CPU es %f%% [Ocupado: %f, Total: %f]\n", cpuUsage, totalTicks-idleTicks, totalTicks)

    stat, err := ReadStat("/proc/stat")
    if err != nil {
        log.Fatal("stat read fail")
    }
    fmt.Printf("Numero de Total de Procesos: %d\n",stat.Processes)
    fmt.Printf("Numero de Procesos en Ejecucion: %d\n",stat.ProcsRunning)
    fmt.Printf("Numero de Procesos Bloqueados: %d\n",stat.ProcsBlocked)

    TestMemInfo()
   
    
}

func TestMemInfo() {
	read, err := ReadMemInfo("/proc/meminfo")
	if err != nil {
		log.Fatal("meminfo read fail: %s",err)
    }
    MemoriaTotal := read.MemTotal/1000
    MemoriaConsumida := ((read.MemTotal-read.MemAvailable)/1000)
    fmt.Printf("Memoria Total del Servidor: %d MB \n",  MemoriaTotal)
    fmt.Printf("Memoria Total Consumida: %d MB\n", MemoriaConsumida)
	fmt.Printf("Porcentaje de Consumo de RAM: %f %%\n", ((float64(MemoriaConsumida)*100)/float64(MemoriaTotal) ))
}
