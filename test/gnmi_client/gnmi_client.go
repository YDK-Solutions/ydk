/*
 * ------------------------------------------------------------------
 * gnmi_client.go
 *
 * May, 2017, Weihuang Fu
 *
 * Copyright (c) 2017-2018 by cisco Systems, Inc.
 * All rights reserved.
 * ------------------------------------------------------------------
 */

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	gnmi "proto/gnmi"
	"runtime"
	"strings"
	"time"
    "bytes"
    "errors"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
)

// Default server name used in EMS certificate (Common Name)
const EMSDefServerName = "ems.cisco.com"
const (
	mega = float64(1000 * 1000)
)

var (
	tls = flag.Bool("tls", false,
		"Connection uses TLS if true, else plain TCP")
	caFile = flag.String("ca_file", "../testdata/ca.pem",
		"The file containning the CA root cert file")
	userName = flag.String("username", "root",
		"Username for the client connection")
	password = flag.String("password", "lab",
		"Password of the username for the client connection")
	serverHostOverride = flag.String("server_host_override", EMSDefServerName,
		"The server name use to verify the hostname returned by TLS handshake")
	cliInFile = flag.String("cli_in_file", "",
		"CLI command input file (applicable for CLI config and CLI action)")
	serverAddr = flag.String("server_addr", "127.0.0.1:57400",
		"The server address in the format of host:port")
	oper = flag.String("oper", "get-config",
		"Operation: get-config, merge-config, replace-config, cli-config, delete-config, cli-show, cli-json, telem-sub, telem-unsub, get-models, action-json, action-cli, commit, commit-replace")
	numOper = flag.Int("num_oper", 1,
		"Number of repeated operations")
	operSleep = flag.Int("oper_sleep", 1,
		"Sleep time in seconds between repeated operations")
	operTimeout = flag.Int("oper_timeout", 5,
		"Timeout for reconnection in seconds")
	grpc_summary = flag.Bool("summary", true,
		"whether to show grpc summary")
	reqIdStart = flag.Int64("req_id_start", 1,
		"Starting request identifier")
)

func trace(msg string) func() {
	start := time.Now()
	//log.Printf(" enter PID:%d:%s %s", os.Getpid(), funcName(2), msg)
	fmt.Printf(" enter PID:%d:%s %s\n", os.Getpid(), funcName(2), msg)
	return func() {
		//log.Printf(" exit PID:%d:%s %s (% s)", os.Getpid(), funcName(2), msg, time.Since(start))
		fmt.Printf(" exit PID:%d:%s %s (% s)\n", os.Getpid(), funcName(2), msg, time.Since(start))
	}
}

func funcName(level int) string {
	pc, _, _, _ := runtime.Caller(level)
	return runtime.FuncForPC(pc).Name()
}

func GetClientMetadata() map[string]string {
	return map[string]string{
		"username": "cisco",
		"password": "cisco123",
	}
}

type passCredential int

func (passCredential) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"username": *userName,
		"password": *password,
	}, nil
}

func (passCredential) RequireTransportSecurity() bool {
	return false
}

func calculateRTT(start, end time.Time, numOper, operSleep float64) float64 {
	var dividor float64
	diff := end.Sub(start)
	diff_secs := diff.Seconds()
	dividor = 1.0
	if numOper != 0 {
		dividor = numOper
	}
	/*
	   p := fmt.Println
	   p(start)
	   p(end)
	   p(diff)
	   p(diff_secs)
	   p(dividor)
	*/
	rtt := (float64(diff_secs) - numOper*operSleep) / dividor

	return rtt
}

func sp(width int) {
	for i := 0; i < width; i++ {
		fmt.Printf("  ")
	}
}

type get struct {
	conn gnmi.GNMIClient
}

type set struct {
	conn     gnmi.GNMIClient
	deletes  []*gnmi.Path
	replaces []*gnmi.Update
	updates  []*gnmi.Update
	rsp      *gnmi.SetResponse
	err      error
}

func (s *set) delete(path *gnmi.Path) {
	s.deletes = append(s.deletes, path)
}

func (s *set) replace(path *gnmi.Path, val *gnmi.TypedValue) {
	update := &gnmi.Update{Path: path, Val: val}
	s.replaces = append(s.replaces, update)
}

func (s *set) update(path *gnmi.Path, val *gnmi.TypedValue) {
	update := &gnmi.Update{Path: path, Val: val}
	s.updates = append(s.updates, update)
}

func goPrint(pre string, ss string) {
	cnt := 0
	fmt.Printf("%s:<\n ", pre)
	for i, c := range ss {
		s := string(c)
		fmt.Printf("%v", s)
		if s == "<" {
			cnt++
		} else if s == ">" {
			cnt--
			if cnt == 0 && i < len(ss)-3 {
				//fmt.Print("i",i,"ss",len(ss))
				fmt.Printf("\n")
			}
		}
	}
	fmt.Printf(">\n")
}

func (s *set) send() {
	req := &gnmi.SetRequest{Delete: s.deletes,
		Replace: s.replaces,
		Update:  s.updates}
	//goPrint("SetRequest", fmt.Sprintln(req))
	printSetRequest(req)
	s.rsp, s.err = s.conn.Set(context.Background(), req)
	if s.err != nil {
		log.Fatalf("ERROR: Set: %v", s.err)
	}

	if s.rsp == nil {
		log.Fatalf("ERROR: Set response nil")
	}

	if (s.rsp.Message != nil) && (s.rsp.Message.Code != uint32(codes.OK)) {
		log.Fatalf("ERROR: Set: %v", s.rsp.Message)
	}

	printSetResponse(s.rsp)
	//goPrint("SetResponse", fmt.Sprintln(s.rsp))
}

type queryType int

const (
	YANG_DELETE queryType = iota
	YANG_REPLACE
	YANG_UPDATE
	YANG_GET
	CLI_CFG
	CLI_SHOW
	RPC_HOLD
	RPC_SEND
	GET_CLI
	GET_CFG
	GET_OPER
	GET_ALL
)

var queryTypeDict = map[string]queryType{
	"u": YANG_UPDATE,
	"r": YANG_REPLACE,
	"d": YANG_DELETE,
	"g": YANG_GET,
	"c": CLI_CFG,
	"s": CLI_SHOW,
}

type tuv struct {
	oper queryType
	url  string
	val  string
}

func (s *tuv) getPath() *gnmi.Path {
	path := &gnmi.Path{}
	if s.oper == CLI_CFG || s.oper == CLI_SHOW {
		path.Origin = "CLI"
		return path
	}
    path, _ =  stringToPathElem(s.url)
    fmt.Println("NAME:", path.Elem[0].Name)
	if len(path.Elem) > 0 {
        i := strings.Index(path.Elem[0].Name, ":")
        if i != -1 {
            path.Origin = path.Elem[0].Name[:i]
            path.Elem[0].Name = path.Elem[0].Name[i+1:]
        }
    }
	return path
}

func (s *tuv) getTypedValue() *gnmi.TypedValue {
	if s.oper == CLI_CFG || s.oper == CLI_SHOW {
		if len(s.url) > 0 {
			dat, err := ioutil.ReadFile(s.url)
			if err != nil {
				fmt.Println("Cannot open", s.url)
				return nil
			}
			s.val = string(dat)
		} else {
			fmt.Println("Cannot open, fname is less than 1")
			return nil
		}
		val := &gnmi.TypedValue{Value: &gnmi.TypedValue_AsciiVal{s.val}}
		return val
	}
	val := &gnmi.TypedValue{Value: &gnmi.TypedValue_JsonIetfVal{[]byte(s.val)}}
	return val
}

func (s *tuv) getUpdate() *gnmi.Update {
	path := s.getPath()
    fmt.Println("set path:", path)
	val := s.getTypedValue()
	update := &gnmi.Update{Path: path, Val: val}
	return update
}

type f2tuv struct {
	fname   string
	queries []tuv
}

func (s *f2tuv) load() error {
	raw, err := ioutil.ReadFile(s.fname)
	dat := string(raw)
	if err == nil {
		lines := strings.Split(dat, "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if len(line) == 0 || string(line[0]) == "#" {
				continue
			}
			if len(line) == 1 {
				if line == "{" {
					s.queries = append(s.queries, tuv{oper: RPC_HOLD})
				} else if line == "}" {
					s.queries = append(s.queries, tuv{oper: RPC_SEND})
				}
				continue
			}
			tmp := strings.Split(line, "@")
			if len(tmp) == 2 {
				s.queries = append(s.queries, tuv{oper: queryTypeDict[tmp[0]], url: tmp[1]})
			}
			if len(tmp) > 2 {
				s.queries = append(s.queries, tuv{oper: queryTypeDict[tmp[0]], url: tmp[1], val: tmp[2]})
			}
		}
	}
	return err
}

type sgc interface {
	delete(path *gnmi.Path)
	replace(path *gnmi.Path, val *gnmi.TypedValue)
	update(path *gnmi.Path, val *gnmi.TypedValue)
	send()
}

type tuv2gnmi struct {
	queries []tuv
	gnmirpc sgc
	conn    gnmi.GNMIClient
}

func (s *tuv2gnmi) run() {
	hold := 0 //0-no_hold 1-hold_init 2-hold
	for _, query := range s.queries {
        fmt.Println("Query:", query)
		switch query.oper {
		case RPC_HOLD:
			hold = 1
		case RPC_SEND:
			hold = 0
		case YANG_UPDATE:
			if hold == 1 {
				s.gnmirpc = &set{conn: s.conn}
				hold = 2
			}
			s.gnmirpc.update(query.getPath(), query.getTypedValue())
		case YANG_DELETE:
			if hold == 1 {
				s.gnmirpc = &set{conn: s.conn}
				hold = 2
			}
			s.gnmirpc.delete(query.getPath())
		case YANG_REPLACE:
			if hold == 1 {
				s.gnmirpc = &set{conn: s.conn}
				hold = 2
			}
			s.gnmirpc.replace(query.getPath(), query.getTypedValue())
		case CLI_CFG:
			if hold == 1 {
				s.gnmirpc = &set{conn: s.conn}
				hold = 2
			}
			s.gnmirpc.update(query.getPath(), query.getTypedValue())
		}
		if hold == 0 {
			s.gnmirpc.send()
		}
	}
}

// stringToPathElem converts a path string to  PathElem format
func stringToPathElem(path string) (*gnmi.Path, error) {
	elm := pathToElement(path)

	gpath := &gnmi.Path{}
	for _, p := range elm {
		name, kv, err := getKV(p)
		if err != nil {
			return nil, fmt.Errorf("error parsing path %s: %v", path, err)
		}
		gpath.Elem = append(gpath.Elem, &gnmi.PathElem{
			Name: name,
			Key:  kv,
		})
	}
	return gpath, nil
}

// pathToElement splits path to elements
func pathToElement(s string) []string {
	var elm []string
	var buf bytes.Buffer

	var isQt, isEsc bool

	for _, ch := range s {
		switch {
        case ch == '[' && !isEsc:
			isQt = true
		case ch == ']' && !isEsc:
			isQt = false
		case ch == '\\' && !isEsc && !isQt:
			isEsc = true
			continue
		case ch == '/' && !isEsc && !isQt:
			elm = append(elm, buf.String())
			buf.Reset()
			continue
		}

		buf.WriteRune(ch)
		isEsc = false
	}

	if buf.Len() != 0 {
		elm = append(elm, buf.String())
	}

	if len(elm) > 0 && elm[0] == "" {
		elm = elm[1:]
	}

	return elm
}

// getKV gets key/value pair from the input string
func getKV(in string) (string, map[string]string, error) {
	var isEsc, isQt, isVal bool
	var name, currentKey string
	var buf bytes.Buffer
	keys := map[string]string{}

	for _, ch := range in {
		switch {
		case ch == '[' && !isEsc && !isVal && isQt:
			return "", nil, fmt.Errorf("received an unescaped [ in key of element %s", name)
		case ch == '[' && !isEsc && !isQt:
			isQt = true
			if len(keys) == 0 {
				if buf.Len() == 0 {
					return "", nil, errors.New("received a value when the element name was null")
				}
				name = buf.String()
				buf.Reset()
			}
			continue
		case ch == ']' && !isEsc && !isQt:
			return "", nil, fmt.Errorf("received an unescaped ] when not in a key for element %s", buf.String())
		case ch == ']' && !isEsc:
			isQt = false
			isVal = false
			if err := addKey(keys, name, currentKey, buf.String()); err != nil {
				return "", nil, err
			}
			buf.Reset()
			currentKey = ""
			continue
		case ch == '\\' && !isEsc:
			isEsc = true
			continue
		case ch == '=' && isQt && !isEsc && !isVal:
			currentKey = buf.String()
			buf.Reset()
			isVal = true
			continue
		}

		buf.WriteRune(ch)
		isEsc = false
	}

	if len(keys) == 0 {
		name = buf.String()
	}

	if len(keys) != 0 && buf.Len() != 0 {
		// In this case, we have trailing garbage following the key.
		return "", nil, fmt.Errorf("trailing garbage following keys in element %s, got: %v", name, buf.String())
	}

	if strings.Contains(name, " ") {
		return "", nil, fmt.Errorf("invalid space character included in element name '%s'", name)
	}

	return name, keys, nil
}

// addKey adds key k with value v to the key map
func addKey(keys map[string]string, e, k, v string) error {
	switch {
	case strings.Contains(k, " "):
		return fmt.Errorf("received an invalid space in element %s key name '%s'", e, k)
	case e == "":
		return fmt.Errorf("received null element value with key and value %s=%s", k, v)
	case k == "":
		return fmt.Errorf("received null key name for element %s", e)
	case v == "":
		return fmt.Errorf("received null value for key %s of element %s", k, e)
	}
	keys[k] = v
	return nil
}
//fname -> [(oper, path, value)] -> set/get -> send

func gnmiMix(client gnmi.GNMIClient, fname string) {
	fh := &f2tuv{fname: fname}
	err := fh.load()
	if err != nil {
		fmt.Println("Error to open file", err)
		return
	}
	rpc := &tuv2gnmi{conn: client, queries: fh.queries}
	rpc.run()
}

func main() {
	//defer trace(" ")()
	flag.Parse()
	var opts []grpc.DialOption
	var cred passCredential
	var minRTT, maxRTT float64
	var reqStart, reqEnd time.Time
	var totalByteTransfer int64

	minRTT = 0.0
	maxRTT = 0.0
	totalByteTransfer = 0

	args := os.Args[1:]
	if args[0] == "get" || args[0] == "set" || args[0] == "cap" ||
		args[0] == "getCli" || args[0] == "getCfg" || args[0] == "getOper" || args[0] == "getAll" {
		if len(args) >= 2 {
			*serverAddr = args[1]
			fmt.Printf("RPC to %v\n", *serverAddr)
		}
		if len(args) >= 3 {
			if args[0] == "cap" && args[2] == "true" {
				*tls = true
				*caFile = args[3]
			} else {
				*cliInFile = args[2]
				fmt.Printf("Input: %v\n", *cliInFile)
			}
		}
		if len(args) > 4 {
			if args[3] == "true" {
				*tls = true
				*caFile = args[4]

			}
		}
	}

	if *tls {
		fmt.Printf("Enabled TLS for Router Application Client\n")
		var sn string

		if *serverHostOverride != "" {
			sn = *serverHostOverride
		}

		var creds credentials.TransportCredentials
		if *caFile != "" {
			var err error
			creds, err = credentials.NewClientTLSFromFile(*caFile, sn)
			if err != nil {
				log.Fatalf("ERROR: Failed to create TLS credentials:", err)
			}
		} else {
			creds = credentials.NewClientTLSFromCert(nil, sn)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		fmt.Printf("TLS is disabled. Using insecure connection for Router Application Client\n")
		opts = append(opts, grpc.WithInsecure())
	}
	opts = append(opts, grpc.WithPerRPCCredentials(cred))
	opts = append(opts, grpc.WithTimeout(time.Second*time.Duration(*operTimeout)))
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("ERROR: Failed to connect to server:", err)
                return
	}
	defer conn.Close()

	gnmiClient := gnmi.NewGNMIClient(conn)
	reqId := *reqIdStart
	// Set starting reqId to Pid if it's the same as the default value
	// We do this to support end-to-end tracing
	if reqId == 1 {
		reqId = int64(os.Getpid())
	}
	//getCli
	//getCfg
	//getOper
	//getAll
	if args[0] == "get" {
		fname := *cliInFile
		if len(fname) >= 3 && fname[len(fname)-3:] == "cli" {
			gnmiGetCli(gnmiClient)
		} else {
			gnmiGet(gnmiClient, GET_CFG)
		}
	} else if args[0] == "getCli" {
		gnmiGetCli(gnmiClient)
	} else if args[0] == "getCfg" {
		gnmiGet(gnmiClient, GET_CFG)
	} else if args[0] == "getOper" {
		gnmiGet(gnmiClient, GET_OPER)
	} else if args[0] == "getAll" {
		gnmiGet(gnmiClient, GET_ALL)
	} else if args[0] == "set" {
		fname := *cliInFile
		if len(fname) >= 3 && fname[len(fname)-3:] == "cli" {
			gnmiSet(gnmiClient)
		} else {
			gnmiMix(gnmiClient, args[2])
		}
	} else if args[0] == "cap" {
		gnmiCap(gnmiClient)
	} else if args[0] == "mix" {
		gnmiMix(gnmiClient, args[2])
	}

	return

	totalStart := time.Now()
	for ii := 0; ii < *numOper; ii++ {
		reqStart = time.Now()
		if strings.EqualFold(*oper, "cap") {
			gnmiCap(gnmiClient)
		} else if strings.EqualFold(*oper, "get") {
			gnmiGet(gnmiClient, GET_CFG)
		} else if strings.EqualFold(*oper, "set") {
			gnmiSet(gnmiClient)
		} else {
			fmt.Printf("Error: unsupported operation %s\n", *oper)
		}

		reqEnd = time.Now()
		rtt := calculateRTT(reqStart, reqEnd, 1.0, 0.0)
		if ii == 0 {
			minRTT = rtt
			maxRTT = rtt
		} else {
			if minRTT > rtt {
				minRTT = rtt
			}
			if maxRTT < rtt {
				maxRTT = rtt
			}
		}
		time.Sleep(time.Second * time.Duration(*operSleep))
		reqId++
	}
	totalEnd := time.Now()

	rtt := minRTT
	totalRtt := rtt

	if *numOper > 1 {
		/* remove sleep time in calculateRTT */
		rtt = calculateRTT(totalStart, totalEnd, float64(*numOper), float64(*operSleep))
		totalRtt = rtt * float64(*numOper)
	}

	if *grpc_summary {
		fmt.Printf("\n----------------- gRPC Summary ----------------------\n\n")
		fmt.Printf("Operation: %s\n", *oper)
		fmt.Printf("Number of iterations: %d\n", *numOper)
		fmt.Printf("Total bytes transferred: %d\n", totalByteTransfer)
		if totalRtt != 0.0 {
			totalBitsTransfer := totalByteTransfer * 8
			fmt.Printf("Number of bytes per second: %d\n", int64(float64(totalByteTransfer)/totalRtt))
			fmt.Printf("Round trip throughputs Mbps: %f\n", float64(totalBitsTransfer)/(totalRtt*mega))
		}
		fmt.Printf("Ave elapsed time in seconds: %f\n", rtt)
		fmt.Printf("Min elapsed time in seconds: %f\n", minRTT)
		fmt.Printf("Max elapsed time in seconds: %f\n", maxRTT)
		fmt.Printf("\n--------------- End gRPC Summary ---------------------\n\n")
	}
}

func gnmiCap(client gnmi.GNMIClient) {
	var req gnmi.CapabilityRequest
	fmt.Println("=== CapabilityRequest ===")
	reply, err := client.Capabilities(context.Background(), &req)
	if err != nil {
		log.Fatalf("ERROR: Cap: %v", err)
	}
	fmt.Println("=== CapabilityResponse ===")
	fmt.Println("<")
	for _, el := range reply.SupportedModels {
		fmt.Printf("  supported_models: <\n    name: %s\n    organization: %s\n    version: %s\n  >\n",
			el.Name, el.Organization, el.Version)
	}
	for _, el := range reply.SupportedEncodings {
		fmt.Printf("  supported_encodings: %s\n", el)
	}
	fmt.Printf("  gNMI_version: %s\n", reply.GNMIVersion)
	fmt.Println(">")
}

func gnmiSet(client gnmi.GNMIClient) {
	var req gnmi.SetRequest
	updates := []*gnmi.Update{}
	//cli-config
	if len(*cliInFile) != 0 {
		update := &gnmi.Update{}
		dat, err := ioutil.ReadFile(*cliInFile)
		if err != nil {
			fmt.Printf("can't open CLI file %s\n", *cliInFile)
			return
		}
		update.Val = &gnmi.TypedValue{Value: &gnmi.TypedValue_AsciiVal{string(dat)}}
		updates = append(updates, update)
	}
	//yang-config
	req.Update = updates

	fmt.Println("=== setRequest ===")
	fmt.Println("setRequest: <")
	if req.GetPrefix() != nil {
		fmt.Println("  prefix: <\n")
		fmt.Println("    element: ")
		fmt.Println("    origin: ")
		fmt.Printf("  >\n>\n")
	}
	if req.GetDelete() != nil {
		fmt.Println("in construction")
	}
	if req.GetReplace() != nil {
		fmt.Println("in construction")
	}
	if req.GetUpdate() != nil {
		for _, update := range req.GetUpdate() {
			fmt.Printf("  update: <\n    val: <\n      ascii_val:\n%s    >\n  >\n>\n", update.GetVal().GetAsciiVal())
		}
	}

	reply, err := client.Set(context.Background(), &req)

	fmt.Println("=== setResponse ===")
	if err != nil {
		log.Fatalf("ERROR: Set: %v", err)
	}

	fmt.Println("setResponse: <")
	//updateResult
	for _, e := range reply.GetResponse() {
		fmt.Printf("  response: <\n")
		fmt.Printf("    message: <\n")
		fmt.Printf("      code: %d\n", e.GetMessage().GetCode())
		fmt.Printf("    >\n")
		fmt.Printf("    op: %s\n", e.GetOp().String())
		fmt.Printf("  >\n")
	}

	fmt.Printf("  timestamp: %d\n>\n", reply.GetTimestamp())
}

func gnmiGet(client gnmi.GNMIClient, tp queryType) {
	var req gnmi.GetRequest
	var plist []string
    var requestInput string
	
    if len(*cliInFile) != 0 {
		dat, err := ioutil.ReadFile(*cliInFile)
		if err != nil {
			fmt.Printf("can't open CLI file %s\n", *cliInFile)
			return
		}
        requestInput = strings.TrimSuffix(string(dat),"\n")
		if len(requestInput) == 1 && requestInput == "/" {
				//request from root
				pt := &gnmi.Path{}
                pathElem := new(gnmi.PathElem)
                pathElem.Name = ""
				pt.Elem  =  append(pt.Elem, pathElem)
				req.Path = append(req.Path, pt)
				goto GET_R
		}
		plist = strings.Split(requestInput, "\n")
	}

	for _, p := range plist {
		pt := &gnmi.Path{}
		if len(p) < 1 {
			break
		}
        pt, _ = stringToPathElem(p)
        if pt == nil {
            fmt.Errorf("Error converting path to path elem")
            break
        }
        i := strings.Index(pt.Elem[0].Name, ":")
        if i != -1 {
            pt.Origin = pt.Elem[0].Name[:i]
            pt.Elem[0].Name = pt.Elem[0].Name[i+1:]
        }
		req.Path = append(req.Path, pt)
	}
GET_R:
	req.Encoding = gnmi.Encoding_JSON_IETF
	if tp == GET_CFG {
		req.Type = gnmi.GetRequest_CONFIG
	} else if tp == GET_OPER {
		req.Type = gnmi.GetRequest_OPERATIONAL
	} else if tp == GET_ALL {
		req.Type = gnmi.GetRequest_ALL
	}

	var pre gnmi.Path
	if len(req.Path) > 1 {
		org := req.Path[0].Origin
		for _, pth := range req.Path {
			if pth.Origin != org {
				goto GET_NEXT
			}
		}
        
        idx := 0
		for i, elm := range req.Path[0].GetElem() {
			for j, pth := range req.Path {
				if j == 0 {
					continue
				}
				e := pth.GetElem()
				if e == nil || len(e) < i+1 {
					goto GET_NEXT
				}
				if e[i].Name != elm.Name {
					if i == 0 {
						goto GET_NEXT
					} else {
						idx = i
						goto GET_PRE
					}
				}
			}
		}
		idx = len(req.Path[0].GetElem())

	GET_PRE:
		pre.Elem = req.Path[0].Elem[:idx]
		pre.Origin = req.Path[0].Origin
		req.Prefix = &pre
		for i, pth := range req.Path {
			req.Path[i].Elem = pth.Elem[idx:]
			req.Path[i].Origin = ""
		}
	}
GET_NEXT:
	fmt.Println("=== getRequest ===")
	printGetRequest(req)

	reply, err := client.Get(context.Background(), &req)
	fmt.Println("=== getResponse ===")
	if err != nil {
		log.Fatalf("ERROR: Get: %v", err)
	}
	printGetResponse(reply, true)
}

func gnmiGetCli(client gnmi.GNMIClient) {
	var req gnmi.GetRequest
	pth := &gnmi.Path{}
	//cli-show
	if len(*cliInFile) != 0 {
		dat, err := ioutil.ReadFile(*cliInFile)
		if err != nil {
			fmt.Printf("can't open CLI file %s\n", *cliInFile)
			return
		}
		tmp := strings.Split(string(dat), "\n")
		for _, e := range tmp {
			if len(e) > 0 {
                pathElem := new(gnmi.PathElem)
                pathElem.Name = e
				pth.Elem = append(pth.Elem, pathElem)
			}
		}
	}

	req.Path = []*gnmi.Path{}
	req.Path = append(req.Path, pth)
	req.Type = gnmi.GetRequest_ALL
	req.Encoding = gnmi.Encoding_ASCII

	fmt.Println("=== getRequest ===")
	printGetRequest(req)

	reply, err := client.Get(context.Background(), &req)

	fmt.Println("=== getResponse ===")
	if err != nil {
		log.Fatalf("ERROR: Get: %v", err)
	}
	printGetResponse(reply, false)
}

func printGetRequest(req gnmi.GetRequest) {
	fmt.Println("getRequest: <")
	if req.GetPrefix() != nil {
		pre := req.GetPrefix()
		fmt.Println("  prefix: <")
        for _, e := range pre.GetElem() {
            fmt.Println("    elem: <")
            fmt.Printf("        name: \"%s\"\n", e.Name)
            if e.Key != nil {
                for k, v := range e.Key {
                    fmt.Printf("         key: <\n")
                    fmt.Printf("            key: \"%s\"\n", k)
                    fmt.Printf("            value: %s\n", v)
                    fmt.Printf("         >\n")//key
                }
            }
            fmt.Println("    >")
        }
		if len(pre.GetOrigin()) > 0 {
			fmt.Printf("    origin: %s\n", pre.GetOrigin())
		}
		fmt.Println("  >")
	}
	for _, pth := range req.GetPath() {
		fmt.Println("  path: <")
        for _, e := range pth.GetElem() {
            fmt.Println("    elem: <")
            fmt.Printf("        name: \"%s\"\n", e.Name)
            if e.Key != nil {
                for k, v := range e.Key {
                    fmt.Printf("         key: <\n")
                    fmt.Printf("            key: \"%s\"\n", k)
                    fmt.Printf("            value: %s\n", v)
                    fmt.Printf("         >\n")//key
                }
            }
            fmt.Println("    >")
        }
		if len(pth.GetOrigin()) > 0 {
			fmt.Printf("    origin: %s\n", pth.GetOrigin())
		}
		fmt.Printf("  >\n")
	}
	fmt.Printf("  type: %s\n", req.GetType().String())
	fmt.Printf("  encoding: %s\n", req.GetEncoding().String())
	fmt.Printf(">\n")
}

func printSetRequest(req *gnmi.SetRequest) {
	fmt.Println("setRequest: <")
	if req.GetPrefix() != nil {
		pre := req.GetPrefix()
		fmt.Println("  prefix: <")
		for _, e := range pre.GetElem() {
	        fmt.Println("    elem: <")
            fmt.Printf("        name: \"%s\"\n", e.Name)
            if e.Key != nil {
                for k, v := range e.Key {
                    fmt.Printf("         key: <\n")
                    fmt.Printf("            key: \"%s\"\n", k)
                    fmt.Printf("            value: %s\n", v)
                    fmt.Printf("         >\n")//key
                }
            }
            fmt.Println("    >")
	
        }
		if len(pre.GetOrigin()) > 0 {
			fmt.Printf("    origin: %s\n", pre.GetOrigin())
		}
		fmt.Println("  >")
	}
	for _, pth := range req.GetDelete() {
		fmt.Println("  delete: <")
		fmt.Println("    path: <")
		for _, e := range pth.GetElem() {
		    fmt.Println("    elem: <")
            fmt.Printf("        name: \"%s\"\n", e.Name)
            if e.Key != nil {
                for k, v := range e.Key {
                    fmt.Printf("         key: <\n")
                    fmt.Printf("            key: \"%s\"\n", k)
                    fmt.Printf("            value: %s\n", v)
                    fmt.Printf("         >\n")//key
                }
            }
            fmt.Println("    >")

        }
		if len(pth.GetOrigin()) > 0 {
			fmt.Printf("      origin: %s\n", pth.GetOrigin())
		}
		fmt.Printf("    >\n")
		fmt.Printf("  >\n")
	}
	for i := 0; i < 2; i++ {
		var updates []*gnmi.Update
		if i == 0 {
			updates = req.GetReplace()
		} else {
			updates = req.GetUpdate()
		}
		for _, update := range updates {
			if i == 0 {
				fmt.Printf("  replace: <\n")
			} else {
				fmt.Printf("  update: <\n")
			}
			fmt.Println("    path: <")
			for _, e := range update.GetPath().GetElem() {
                fmt.Println("    elem: <")
                fmt.Printf("        name: \"%s\"\n", e.Name)
                if e.Key != nil {
                    for k, v := range e.Key {
                        fmt.Printf("         key: <\n")
                        fmt.Printf("            key: \"%s\"\n", k)
                        fmt.Printf("            value: %s\n", v)
                        fmt.Printf("         >\n")//key
                    }
                }
                fmt.Println("    >")

			}
			if len(update.GetPath().GetOrigin()) > 0 {
				fmt.Printf("      origin: %s\n", update.GetPath().GetOrigin())
			}
			fmt.Printf("    >\n")
			switch update.GetVal().GetValue().(type) {
			case *gnmi.TypedValue_JsonIetfVal:
				fmt.Printf("    val: <\n      json_ietf_val: %s\n", string(update.GetVal().GetJsonIetfVal()))
			case *gnmi.TypedValue_AsciiVal:
				fmt.Printf("    val: <\n      ascii_val: %s", string(update.GetVal().GetAsciiVal()))
			}
			fmt.Printf("    >\n")
			fmt.Printf("  >\n")
		}
	}

	fmt.Printf(">\n")
}

func printGetResponse(reply *gnmi.GetResponse, yang bool) {
	fmt.Printf("getResponse: <\n")
	for _, notif := range reply.GetNotification() {
		fmt.Printf("  notification: <\n")
		fmt.Printf("    timestamp: %d\n", notif.GetTimestamp())
		if notif.GetPrefix() != nil {
			pre := notif.GetPrefix()
			fmt.Println("    prefix: <")
			for _, e := range pre.GetElem() {
                fmt.Printf("        elem: <\n")
                fmt.Printf("            name: \"%s\"\n", e.Name)
                if e.Key != nil {
                    for k,v := range e.Key {
                        fmt.Printf("            key: <\n")
                        fmt.Printf("               key: \"%s\"\n", k)
                        fmt.Printf("               value: %s\n", v)
                        fmt.Printf("            >\n")//key

                    }
                }
                fmt.Printf("        >\n")//elem

			}
			if len(pre.GetOrigin()) > 0 {
				fmt.Printf("      origin: %s\n", pre.GetOrigin())
			}
			fmt.Println("    >")
		}
        for _, update := range notif.GetUpdate() {
			fmt.Printf("    update: <\n")
			fmt.Printf("      path: <\n")
			for _, e := range update.GetPath().GetElem() {
				fmt.Printf("        elem: <\n")
                fmt.Printf("            name: \"%s\"\n", e.Name)
                if e.Key != nil {
                    for k,v := range e.Key {
                        fmt.Printf("            key: <\n")
                        fmt.Printf("               key: \"%s\"\n", k)
                        fmt.Printf("               value: %s\n", v)
                        fmt.Printf("            >\n")//key

                    }
                }
                fmt.Printf("        >\n")//elem
			}
			if len(update.GetPath().GetOrigin()) > 0 {
				fmt.Printf("        origin: %s\n", update.GetPath().GetOrigin())
			}
			fmt.Printf("      >\n")
			if yang {
				fmt.Printf("      val: <\n        json_ietf_val: %s\n", string(update.GetVal().GetJsonIetfVal()))
			} else {
				fmt.Printf("      val: <\n        ascii_val: %s", string(update.GetVal().GetAsciiVal()))
			}
			fmt.Println("      >") //val
			fmt.Println("    >")   //update
		}
    }
	fmt.Println(">")
}

func printSetResponse(reply *gnmi.SetResponse) {
	fmt.Printf("setResponse: <\n")
	for _, resp := range reply.GetResponse() {
		fmt.Printf("  response: <\n")
		fmt.Printf("    path: <\n")
		if resp.GetPath() != nil {
			pth := resp.GetPath()
            for _, e := range pth.GetElem() {
                fmt.Printf("        elem: <\n")
                fmt.Printf("            name: \"%s\"\n", e.Name)
                if e.Key != nil {
                    for k,v := range e.Key {
                        fmt.Printf("            key: <\n")
                        fmt.Printf("               key: \"%s\"\n", k)
                        fmt.Printf("               value: %s\n", v)
                        fmt.Printf("            >\n")//key

                    }
                }
                fmt.Printf("        >\n")//elem

			}
			if len(pth.GetOrigin()) > 0 {
				fmt.Printf("      origin: %s\n", pth.GetOrigin())
			}
			fmt.Println("    >")
		}
		fmt.Printf("    op: %s\n", resp.GetOp().String())
		fmt.Println("  >")
	}
	fmt.Printf("  timestamp: %d\n", reply.GetTimestamp())
	fmt.Println(">")
}

