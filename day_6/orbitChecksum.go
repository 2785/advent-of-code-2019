package main

import "strings"

import "fmt"

import "io/ioutil"

type nodeOrbits struct {
	directOrbit   string
	indirectOrbit map[string]struct{}
}

type nodes map[string]nodeOrbits

func parseOrbitInputIntoNodes(orbitData string) nodes {

	orbits := strings.Fields(orbitData)
	nodeMap := make(nodes, len(orbits)+1)
	for _, orbit := range orbits {
		split := strings.Split(orbit, ")")
		nodeMap[split[1]] = nodeOrbits{directOrbit: split[0], indirectOrbit: make(map[string]struct{})}
	}
	return nodeMap
}

func getOrbitCount(nodes nodes) (nodes, int, int) {
	shouldGoOn := true
	for shouldGoOn {
		// Default it to false, change it back to true if modification performed in loop
		shouldGoOn = false
		for i, node := range nodes {
			if node.directOrbit != "COM" {
				parentNode := nodes[node.directOrbit]
				if len(node.indirectOrbit) != len(parentNode.indirectOrbit)+1 {
					for k, v := range parentNode.indirectOrbit {
						nodes[i].indirectOrbit[k] = v
					}
					nodes[i].indirectOrbit[parentNode.directOrbit] = struct{}{}
					shouldGoOn = true
				}
			}
		}

	}
	directOrbitCount, indirectOrbitCount := len(nodes), 0
	for _, node := range nodes {
		indirectOrbitCount += len(node.indirectOrbit)
	}
	return nodes, directOrbitCount, indirectOrbitCount
}

func main() {
	// testInput := "COM)B\nB)C\nC)D\nD)E\nE)F\nB)G\nG)H\nD)I\nE)J\nJ)K\nK)L"
	// nodes := parseOrbitInputIntoNodes(testInput)
	// directOrbitCount, indirectOrbitCount := getOrbitCount(nodes)
	// fmt.Print(nodes, "\nDirect Orbit Count: ", directOrbitCount, "   Indirect Orbit Count: ", indirectOrbitCount, "\nTotal Orbit Count: ", indirectOrbitCount+directOrbitCount)
	file, e := ioutil.ReadFile("orbit_data.txt")
	if e != nil {
		fmt.Println("Error: ", e.Error())
	}
	nodes := parseOrbitInputIntoNodes(string(file))
	nodes, directOrbitCount, indirectOrbitCount := getOrbitCount(nodes)
	fmt.Print("\nDirect Orbit Count: ", directOrbitCount, "   Indirect Orbit Count: ", indirectOrbitCount, "\nTotal Orbit Count: ", indirectOrbitCount+directOrbitCount)
	you := nodes["YOU"].indirectOrbit
	san := nodes["SAN"].indirectOrbit
	youSanTotalIndirectOrbitCount := len(you) + len(san)
	commonObject, totalPath := "", youSanTotalIndirectOrbitCount
	for k := range you {
		_, ok := san[k]
		if ok {
			d := youSanTotalIndirectOrbitCount - 2*len(nodes[k].indirectOrbit)
			if d < totalPath {
				totalPath = d
				commonObject = k
			}
		}
	}
	fmt.Print("\nCommon object in the tree: ", commonObject, "\nTotal Path: ", totalPath-2)
}
