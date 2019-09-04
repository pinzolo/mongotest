package mongotest

// DocData is document data
//   key: field name
//   value: field value
type DocData map[string]interface{}

// StringValue returns value of given key as string.
func (d DocData) StringValue(key string) (string, bool) {
	v, ok := d[key]
	if !ok {
		return "", false
	}
	s, ok := v.(string)
	if !ok || s == "" {
		return "", false
	}
	return s, true
}

// CollectionData is collection data
//   key: document ID
//   value: document data (exclude ID)
type CollectionData map[string]DocData

// DataSet is collection of collection data
//   key: collection name
//   value: collection data
type DataSet map[string]CollectionData
