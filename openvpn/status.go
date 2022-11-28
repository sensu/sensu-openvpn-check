package openvpn

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	clientListHdr       = "OpenVPN CLIENT LIST"
	routingTableHdr     = "ROUTING TABLE"
	globalStatsHdr      = "GLOBAL STATS"
	updatedField        = "Updated"
	maxQueueLengthField = "Max bcast/mcast queue length"
)

const (
	clientListHeaders = 1 << iota
	routingTableHeaders
	globalStatsHeaders
)

// GlobalStats struct to store Global stats
type GlobalStats struct {
	MaxBcastMcastQueueLen int
}

// Status struct to store the status
type Status struct {
	ClientCount  uint
	RouteCount   uint
	GlobalStats  GlobalStats
	LastModified time.Time
	IsUp         bool
}

type ParseError struct {
	s string
}

func (e *ParseError) Error() string {
	return e.s
}

func NewParseError(s string) *ParseError {
	return &ParseError{s}
}

var (
	clientListHeaderColumns = [5]string{
		"Common Name",
		"Real Address",
		"Bytes Received",
		"Bytes Sent",
		"Connected Since",
	}

	routingTableHeaderColumns = [4]string{
		"Virtual Address",
		"Common Name",
		"Real Address",
		"Last Ref",
	}
)

func checkHeaders(headers []string) int {
	if checkClientListHeaders(headers) {
		return clientListHeaders
	} else if checkRoutingTableHeaders(headers) {
		return routingTableHeaders
	} else {
		return 0
	}
}

func checkClientListHeaders(headers []string) bool {
	for i, v := range headers {
		if v != clientListHeaderColumns[i] {
			return false
		}
	}
	return true
}

func checkRoutingTableHeaders(headers []string) bool {
	for i, v := range headers {
		if v != routingTableHeaderColumns[i] {
			return false
		}
	}
	return true
}

// ParseFile parses OpenVPN Status file ad returns a Status struct
func ParseFile(file string) (*Status, error) {
	conn, err := os.Open(file)
	defer func() { _ = conn.Close() }()
	if err != nil {
		return &Status{IsUp: false}, err
	}

	info, err := conn.Stat()
	if err != nil {
		return &Status{IsUp: false}, err
	}
	lastModified := info.ModTime()

	reader := bufio.NewReader(conn)
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)
	var clientCount uint
	var routeCount uint
	var maxBcastMcastQueueLen int
	nextFieldType := 0
	isEmpty := true
	for scanner.Scan() {
		isEmpty = false
		fields := strings.Split(scanner.Text(), ",")
		if fields[0] == "END" && len(fields) == 1 {
			// Stats footer.
		} else if fields[0] == clientListHdr {
			// Header
		} else if fields[0] == routingTableHdr {
			// Routing table header
		} else if fields[0] == globalStatsHdr {
			nextFieldType = globalStatsHeaders
		} else if fields[0] == updatedField && len(fields) == 2 {
			// Skip, use file update time instead as time format is not guaranteed
		} else if checkHeaders(fields) == clientListHeaders {
			nextFieldType = clientListHeaders
		} else if checkHeaders(fields) == routingTableHeaders {
			nextFieldType = routingTableHeaders
		} else if nextFieldType == clientListHeaders && len(fields) == 5 {
			clientCount++
		} else if nextFieldType == routingTableHeaders && len(fields) == 4 {
			routeCount++
		} else if nextFieldType == globalStatsHeaders && len(fields) == 2 {
			if fields[0] == maxQueueLengthField {
				i, err := strconv.Atoi(fields[1])
				if err == nil {
					maxBcastMcastQueueLen = i
				}
			}
		} else {
			return &Status{IsUp: false}, NewParseError("Unable to Parse Status file")
		}
	}
	if isEmpty {
		return &Status{IsUp: false}, NewParseError("Status File is empty")
	}
	return &Status{
		ClientCount:  clientCount,
		RouteCount:   routeCount,
		GlobalStats:  GlobalStats{maxBcastMcastQueueLen},
		LastModified: lastModified,
		IsUp:         true,
	}, nil
}
