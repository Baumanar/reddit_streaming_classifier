package cli

import (
	"fmt"
	"strconv"
	"strings"
)

// Broker contains the host name and the port
type Broker struct {
	host string
	port int
}

// String returns the broker as a string
func (b Broker) String() string {
	return fmt.Sprintf("%v:%v", b.host, b.port)
}

// SetBroker sets the broker
func SetBroker(i []string) Broker {

	b := Broker{
		host: "localhost",
		port: 9092}

	// Try to get the broker details from the commandline arg
	// The validation and error checking here is rudimentary.
	if len(i) == 2 {
		// hard coded to expect a single user-provided arg (the first arg is the program name)
		a := i[1]
		// Check that there's a colon (i.e. we're hopefully going to get host:port)
		if p := strings.Split(a, ":"); len(p) == 2 {
			// Check that the port is an integer
			if _, err := strconv.Atoi(p[1]); err == nil {
				b.port, _ = strconv.Atoi(p[1])
				b.host = p[0]
			} else {
				fmt.Printf("(%v is not an integer, so ignoring provided value and defaulting to %v)\n", p[1], b)
			}
		} else {
			fmt.Printf("(Commandline value %v doesn't look like a host:port, so defaulting to %v)\n", a, b)
		}
	} else {
		fmt.Printf("(A single commandline argument should be used to specify the broker. Defaulting to %v)\n", b)
	}
	return b

}
