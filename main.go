// proto-proxyman generates a Proxyman Protobuf configuration file from a
// Protobuf FileDescriptorSet.
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

// TODO: add config for passing a base-url?
func main() {
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Printf("error reading descriptor set: %v", err)
		os.Exit(1)
	}

	req := &descriptorpb.FileDescriptorSet{}
	if proto.Unmarshal(input, req); err != nil {
		log.Printf("error unmarshalling descriptor set: %v", err)
		os.Exit(1)
	}

	methods := []map[string]interface{}{}

	for _, file := range req.File {
		for _, svc := range file.Service {
			for _, mtd := range svc.Method {
				methods = append(methods, config(
					*file.Package,
					*svc.Name,
					*mtd.Name,
					*mtd.InputType,
					*mtd.OutputType,
				))
			}
		}
	}

	b, err := json.Marshal(methods)
	if err != nil {
		log.Printf("error: %v", err)
		os.Exit(1)
	}

	os.Stdout.Write(b)
}

// config builds a config map that Proxyman understands, assuming that:
// * methods maps to a consistent `*/package.service/method` URL
// * method request & response message are all defined in the same file
//   - not sure if this is actually a constraint, but this doesn't do any
//     resolving of messages whatsoever, it just assumes that they are defined
//     fully qualified from the root, and removes the `.` prefix.
func config(pkg, svc, mtd, in, out string) map[string]interface{} {
	return map[string]interface{}{
		"method": map[string]interface{}{
			"exact": []map[string]interface{}{
				map[string]interface{}{"name": "POST"},
			},
		},
		"isEnabled":        true,
		"payloadType":      1, // "Single Message"
		"url":              fmt.Sprintf("*/%v.%v/%v", pkg, svc, mtd),
		"isIncludingPaths": false,         // "Include all subpaths of this URL"
		"regex":            "useWildcard", // "Simple wildcard * and ? are supported"
		"id":               uuid.New().String(),
		"nameForResponse":  strings.TrimPrefix(out, "."),
		"name":             strings.TrimPrefix(in, "."),
	}
}
