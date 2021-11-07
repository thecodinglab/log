package fields

import "fmt"

type KV map[string]interface{}

// From creates a new KV map from the given list of arguments.
func From(pairs ...interface{}) KV {
	if len(pairs)%2 != 0 {
		panic("invalid number of arguments: missing value for key or vice versa")
	}

	if len(pairs) == 0 {
		return nil
	}

	fields := KV{}

	for i := 0; i < len(pairs); i += 2 {
		key := fmt.Sprint(pairs[i])
		val := pairs[i+1]

		fields[key] = val
	}

	return fields
}

func (kv KV) Clone() KV {
	dst := KV{}
	for k, v := range kv {
		if o, ok := v.(KV); ok {
			v = o.Clone()
		}

		dst[k] = v
	}
	return dst
}

// Append merges the given KV into the current instance.
func (kv KV) Append(o KV) {
	for k, v := range o {
		if s, ok := v.(KV); ok {
			if d, ok := kv[k].(KV); ok {
				d.Append(s)
				continue
			}
		}

		kv[k] = v
	}
}

// Merge merges the two given KV into a new KV instance.
func Merge(a, b KV) KV {
	dst := a.Clone()
	dst.Append(b)
	return dst
}
