package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

var total = 0
var checked = 0
var totalChecked = 0

func main() {
start:
	fmt.Print("Ip range to scan (eg: 85.0.0.0 or 85.0.0.0-85.0.1.0): ")

	ipRangeScanner := bufio.NewScanner(os.Stdin)
	ipRangeScanner.Scan()

	ipRange := ipRangeScanner.Text()

	fmt.Print("Port range to scan (eg: 80 or 80-444): ")

	portRangeScanner := bufio.NewScanner(os.Stdin)
	portRangeScanner.Scan()

	portRange := portRangeScanner.Text()

	var ips []string
	if strings.Contains(ipRange, "-") {
		ips = strings.Split(ipRange, "-")
	} else {
		ips = []string{ipRange, ipRange}
	}

	var ports []string
	if strings.Contains(portRange, "-") {
		ports = strings.Split(portRange, "-")
	} else {
		ports = []string{portRange, portRange}
	}

	minPort, err := strconv.Atoi(ports[0])
	maxPort, err := strconv.Atoi(ports[1])

	if err != nil {
		fmt.Println("Error parsing port range.")
	}

	checkRange(ips[0], ips[1], minPort, maxPort, 5)

	fmt.Print("Exit? (yes or no): ")
	exitScanner := bufio.NewScanner(os.Stdin)
	exitScanner.Scan()

	exitRange := exitScanner.Text()

	if exitRange == "no" || exitRange == "n" {
		goto start
	}
}

func checkRange(min string, max string, minPort int, maxPort int, interval time.Duration) {
	startMilli := time.Now().UnixMilli()

	splitMin := strings.Split(min, ".")
	splitMax := strings.Split(max, ".")

	int1, _ := strconv.Atoi(splitMin[0])
	int2, _ := strconv.Atoi(splitMin[1])
	int3, _ := strconv.Atoi(splitMin[2])
	int4, _ := strconv.Atoi(splitMin[3])

	int1max, _ := strconv.Atoi(splitMax[0])
	int2max, _ := strconv.Atoi(splitMax[1])
	int3max, _ := strconv.Atoi(splitMax[2])
	int4max, _ := strconv.Atoi(splitMax[3])

	total = countIps(int1, int2, int3, int4, int1max, int2max, int3max, int4max, minPort, maxPort)
	fmt.Println("Total ports to check:", total)

	go checkInfo()

	for {
		ipStr := fmt.Sprintf("%d.%d.%d.%d", int1, int2, int3, int4)
		port := minPort

		for {
			checked++
			time.Sleep(time.Microsecond * interval)

			go checkConn(ipStr, strconv.Itoa(port))

			if port == maxPort {
				break
			}

			port++
		}

		if int1 == 255 && int2 == 255 && int3 == 255 && int4 == 255 {
			break
		}

		if int1 == int1max && int2 == int2max && int3 == int3max && int4 == int4max {
			break
		}

		if int2 == 255 && int3 == 255 && int4 == 255 {
			int1++
			int2, int3, int4 = 0, 0, 0
			continue
		}

		if int3 == 255 && int4 == 255 {
			int2++
			int3, int4 = 0, 0
			continue
		}

		if int4 == 255 {
			int3++
			int4 = 0
			continue
		}

		int4++
	}

	msToStart := time.Now().UnixMilli() - startMilli
	time.Sleep(1200 * time.Millisecond)
	fmt.Printf("Range %s to %s. Connection checked to %d ips and ports. In %dms\n", min, max, checked, msToStart)
}

func countIps(int1 int, int2 int, int3 int, int4 int, int1max int, int2max int, int3max int, int4max int, minPort int, maxPort int) int {
	count := 0

	for {
		port := minPort

		for {
			count++

			if port == maxPort {
				break
			}

			port++
		}

		if int1 == 255 && int2 == 255 && int3 == 255 && int4 == 255 {
			break
		}

		if int1 == int1max && int2 == int2max && int3 == int3max && int4 == int4max {
			break
		}

		if int2 == 255 && int3 == 255 && int4 == 255 {
			int1++
			int2, int3, int4 = 0, 0, 0
			continue
		}

		if int3 == 255 && int4 == 255 {
			int2++
			int3, int4 = 0, 0
			continue
		}

		if int4 == 255 {
			int3++
			int4 = 0
			continue
		}

		int4++
	}

	return count
}

func checkInfo() {
	for {
		time.Sleep(time.Second * 1)

		totalChecked += checked
		remaining := total - totalChecked

		if remaining < 0 {
			remaining = 0
		}

		fmt.Printf("Remaining: %d Checked: %d Per second: %d\n", remaining, totalChecked, checked)
		checked = 0

		if remaining == 0 || totalChecked >= total {
			break
		}
	}
}

func checkConn(ip string, port string) {
	conn, err := net.DialTimeout("tcp", ip+":"+port, 500*time.Millisecond)

	if err != nil {
		return
	}

	defer func(conn net.Conn) {
		_ = conn.Close()
	}(conn)

	fmt.Printf("Server found in %s:%s.\n", ip, port)
}
