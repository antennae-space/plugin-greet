// The backend of Greet: JSON-RPC 2.0 over stdin and stdout, one message per line. The IDE starts
// it confined (Landlock and seccomp): it reaches its own folder, its data folder, the project as
// granted, and asks the IDE for the rest. It counts the entries at the root of the project
// through the IDE's fs.list (the grant core:fs.list), for the tool's view and the assistant.
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
)

type message struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Method  string          `json:"method,omitempty"`
	Params  json.RawMessage `json:"params,omitempty"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *rpcError       `json:"error,omitempty"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var (
	writeMu sync.Mutex
	out     = json.NewEncoder(os.Stdout)

	// The calls of the backend to the IDE, waiting for their answer.
	pendingMu sync.Mutex
	pending   = map[string]chan message{}
	next      atomic.Int64
)

func send(m message) {
	m.JSONRPC = "2.0"
	writeMu.Lock()
	defer writeMu.Unlock()
	_ = out.Encode(m)
}

func reply(id json.RawMessage, result any) {
	data, err := json.Marshal(result)
	if err != nil {
		fail(id, err)
		return
	}
	send(message{ID: id, Result: data})
}

func fail(id json.RawMessage, err error) {
	send(message{ID: id, Error: &rpcError{Code: -32000, Message: err.Error()}})
}

// core asks the IDE a method of its core, which it refuses without the grant core:<method>.
func core(method string, params any) (json.RawMessage, error) {
	id := "core-" + strconv.FormatInt(next.Add(1), 10)
	idRaw, _ := json.Marshal(id)
	ch := make(chan message, 1)
	pendingMu.Lock()
	pending[string(idRaw)] = ch
	pendingMu.Unlock()
	p, _ := json.Marshal(map[string]any{"method": method, "params": params})
	send(message{ID: idRaw, Method: "core", Params: p})
	m := <-ch
	if m.Error != nil {
		return nil, fmt.Errorf("%s", m.Error.Message)
	}
	return m.Result, nil
}

// count is the number of entries at the root of the project.
func count() (int, error) {
	raw, err := core("fs.list", map[string]any{"path": ""})
	if err != nil {
		return 0, err
	}
	var entries []json.RawMessage
	if err := json.Unmarshal(raw, &entries); err != nil {
		return 0, err
	}
	return len(entries), nil
}

func main() {
	in := bufio.NewScanner(os.Stdin)
	in.Buffer(make([]byte, 1<<20), 8<<20)
	for in.Scan() {
		var m message
		if json.Unmarshal(in.Bytes(), &m) != nil {
			continue
		}
		if m.Method == "" { // the answer of the IDE to one of our calls
			pendingMu.Lock()
			ch := pending[string(m.ID)]
			delete(pending, string(m.ID))
			pendingMu.Unlock()
			if ch != nil {
				ch <- m
			}
			continue
		}
		go handle(m)
	}
}

func handle(m message) {
	switch m.Method {
	case "initialize":
		reply(m.ID, map[string]any{})
	case "shutdown":
		reply(m.ID, map[string]any{})
		os.Exit(0)
	case "call": // a greet.* method of one of its views
		var c struct{ Method string }
		_ = json.Unmarshal(m.Params, &c)
		switch c.Method {
		case "greet.count":
			n, err := count()
			if err != nil {
				fail(m.ID, err)
				return
			}
			reply(m.ID, n)
		default:
			fail(m.ID, fmt.Errorf("unknown method: %s", c.Method))
		}
	case "tool": // its assistant tool
		n, err := count()
		if err != nil {
			reply(m.ID, map[string]any{"content": err.Error(), "summary": "failed", "error": true})
			return
		}
		reply(m.ID, map[string]any{"content": fmt.Sprintf("%d entries at the root of the project.", n), "summary": fmt.Sprintf("%d entries", n)})
	default:
		fail(m.ID, fmt.Errorf("unknown method: %s", m.Method))
	}
}
