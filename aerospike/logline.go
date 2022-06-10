package aerospike

import (
	"regexp"
	"strings"
)

var regexAddr = regexp.MustCompile(`\d+\.\d+\.\d+\.\d+(:\d+|)`)
var regexNodeID = regexp.MustCompile(`NODE-ID .* CLUSTER-SIZE`)

// TODO: handle other required operations here
func (a *Aerospike) LogLine(line string) (newline string, err error) {

	// check for node-id
	res := regexNodeID.FindAllString(line, -1)
	for _, r := range res {
		rr := strings.Split(r, " ")
		rr[1], err = a.NodeID(rr[1])
		if err != nil {
			continue
		}
		line = strings.ReplaceAll(line, r, strings.Join(rr, " "))
	}
	if len(res) > 0 {
		// this line will not contain anything else, might as well exit now
		return line, nil
	}

	// check for addr
	if !strings.Contains(line, "<><><><><><><><>") {
		var rr string
		res = regexAddr.FindAllString(line, -1)
		for _, r := range res {
			if strings.Contains(r, ":") {
				rr, err = a.Addr(r)
			} else {
				rr, err = a.IP(r)
			}
			if err != nil {
				continue
			}
			line = strings.ReplaceAll(line, r, rr)
		}
	}

	// return
	return line, nil
}
