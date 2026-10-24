// The backend of Greet, on the plugin SDK (antennae.space/plugin-sdk/backend): the IDE starts it
// confined (Landlock and seccomp) and speaks JSON-RPC to it over stdio. It reaches its own
// folder, its data folder, the project as granted, and asks the IDE for the rest: it counts the
// entries at the root of the project through the IDE's fs.list (the grant core:fs.list), for
// the tool's view and the assistant.
package main

import (
	"context"
	"encoding/json"
	"fmt"

	"antennae.space/plugin-sdk/backend"
)

// count is the number of entries at the root of the project.
func count(ctx context.Context) (int, error) {
	var entries []json.RawMessage
	if err := backend.Core(ctx, "fs.list", map[string]any{"path": ""}, &entries); err != nil {
		return 0, err
	}
	return len(entries), nil
}

func main() {
	backend.Serve(backend.Plugin{
		// A greet.* method of one of its views.
		Methods: map[string]backend.Method{
			"greet.count": func(ctx context.Context, _ json.RawMessage) (any, error) { return count(ctx) },
		},
		// Its assistant tool.
		Tools: map[string]backend.Tool{
			"greet_count": func(ctx context.Context, _ json.RawMessage) (backend.Result, error) {
				n, err := count(ctx)
				if err != nil {
					return backend.Result{}, err
				}
				return backend.Result{Content: fmt.Sprintf("%d entries at the root of the project.", n), Summary: fmt.Sprintf("%d entries", n)}, nil
			},
		},
	})
}
