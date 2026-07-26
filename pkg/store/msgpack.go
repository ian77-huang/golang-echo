package store

import "github.com/vmihailenco/msgpack/v5"

func msgpackEncode(data interface{}) ([]byte, error) {
	bVal, err := msgpack.Marshal(&StoreValue{Val: data})
	if err != nil {
		return nil, err
	}
	return bVal, nil
}
func msgpackDecode(data []byte, target interface{}) error {
	storeVal := StoreValue{
		Val: target,
	}
	return msgpack.Unmarshal(data, &storeVal)
}
