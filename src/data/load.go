package data

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/TylerBrock/colorjson"
)

func LoadJson(r io.Reader) (*Root, error) {
	var root = &Root{}

	var r2 bytes.Buffer
	r1 := io.TeeReader(r, &r2)

	var dec = json.NewDecoder(r1)

	if err := dec.Decode(&root.Cells); err != nil {
		var serr *json.SyntaxError
		if errors.As(err, &serr) {
			return nil, fmt.Errorf("json syntax error: %w: string:\n%v", err, r2.String())
		}
		var obj interface{}
		var debugDecoder = json.NewDecoder(&r2)
		debugDecoder.Decode(&obj)
		f := colorjson.NewFormatter()
		f.Indent = 2
		s, _ := f.Marshal(obj)
		return nil, fmt.Errorf("%w - object: %s", err, s)
	}

	return root, nil
}
