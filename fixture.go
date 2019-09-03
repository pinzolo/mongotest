package mongotest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

// docData is document data
//   key: field name
//   value: field value
type docData map[string]interface{}

// collData is collection data
//   key: document ID
//   value: document data (exclude ID)
type collData map[string]docData

// dataSet is collection of document data
//   key: collection name
//   value: collection data
type dataSet map[string]collData

// UseFixtureWithContext apply fixture data to MongoDB with context.Context.
// If multi names are given, fixture data will be merged.(overwriting by after dataset)
func UseFixtureWithContext(ctx context.Context, names ...string) error {
	if err := validateConfig(); err != nil {
		return err
	}
	files, err := toFilePaths(names...)
	if err != nil {
		return err
	}
	ds, err := loadDataSet(files...)
	if err != nil {
		return err
	}
	bs, _ := yaml.Marshal(ds)
	fmt.Println(string(bs))
	for cn, cd := range ds {
		err = resetCollection(ctx, cn, toValues(cd))
		if err != nil {
			return err
		}
	}
	return nil
}

// UseFixture apply fixture data to MongoDB.
// If multi names are given, fixture data will be merged.(overwriting by after dataset)
func UseFixture(names ...string) error {
	return UseFixtureWithContext(context.Background(), names...)
}

func toFilePaths(names ...string) ([]string, error) {
	files := make([]string, len(names))
	for i, name := range names {
		file, err := findFixtureFilePath(name)
		if err != nil {
			return nil, err
		}
		files[i] = file
	}
	return files, nil
}

func findFixtureFilePath(name string) (string, error) {
	dir, base := fixturePath(name)
	fis, err := ioutil.ReadDir(dir)
	if err != nil {
		return "", err
	}
	for _, fi := range fis {
		if fi.Name() == base+filepath.Ext(fi.Name()) && !fi.IsDir() {
			return filepath.Join(dir, fi.Name()), nil
		}
	}
	return "", fmt.Errorf("dataSet %q not found in %s", name, conf.fixtureRootDirAbs)
}

func fixturePath(name string) (dir string, base string) {
	paths := splitDataSetName(name)
	if len(paths) == 1 {
		return conf.fixtureRootDirAbs, name
	}
	dirPaths := []string{conf.fixtureRootDirAbs}
	dirPaths = append(dirPaths, paths[0:len(paths)-1]...)
	return filepath.Join(dirPaths...), paths[len(paths)-1]
}

func splitDataSetName(name string) []string {
	if strings.Contains(name, "/") {
		return strings.Split(name, "/")
	}
	return []string{name}
}

func loadDataSet(files ...string) (dataSet, error) {
	dss, err := toDataSets(files...)
	if err != nil {
		return nil, err
	}
	return mergeDataSet(dss), nil
}

func toDataSets(files ...string) ([]dataSet, error) {
	dss := make([]dataSet, len(files))
	for i, file := range files {
		ds, err := readFixtureFile(file)
		if err != nil {
			return nil, err
		}
		dss[i] = ds
	}
	return dss, nil
}

func readFixtureFile(file string) (dataSet, error) {
	format, err := fixtureFormat(file)
	if err != nil {
		return nil, err
	}
	bs, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	ds := make(dataSet)
	if format == FixtureFormatYAML {
		err = yaml.Unmarshal(bs, &ds)
	} else if format == FixtureFormatJSON {
		err = json.Unmarshal(bs, &ds)
	}
	if err != nil {
		return nil, err
	}
	return ds, nil
}

func fixtureFormat(file string) (FixtureFormatType, error) {
	if conf.fixtureFormat != FixtureFormatAuto {
		return conf.fixtureFormat, nil
	}
	ext := strings.ToLower(filepath.Ext(file))
	switch ext {
	case ".json":
		return FixtureFormatJSON, nil
	case ".yaml", ".yml":
		return FixtureFormatYAML, nil
	default:
		return fixtureFormatUnknown, errors.New("unknown format")
	}
}

func mergeDataSet(dss []dataSet) dataSet {
	merged := make(dataSet)
	for _, ds := range dss {
		for k, v := range ds {
			if coll, ok := merged[k]; ok {
				merged[k] = mergeCollData(coll, v)
			} else {
				merged[k] = v
			}
		}
	}
	return merged
}

func mergeCollData(coll1, coll2 collData) collData {
	merged := make(collData)
	for k, v := range coll1 {
		merged[k] = v
	}
	for k, v := range coll2 {
		if doc, ok := merged[k]; ok {
			merged[k] = mergeDocData(doc, v)
		} else {
			merged[k] = v
		}
	}
	return merged
}

func mergeDocData(doc1, doc2 docData) docData {
	merged := make(docData)
	for k, v := range doc1 {
		merged[k] = v
	}
	for k, v := range doc2 {
		merged[k] = v
	}
	return merged
}

func toValues(coll collData) []interface{} {
	values := make([]interface{}, 0, len(coll))
	for id, doc := range coll {
		newDoc := make(docData)
		for k, v := range doc {
			newDoc[k] = v
		}
		newDoc["_id"] = id
		values = append(values, newDoc)
	}
	return values
}

func resetCollection(ctx context.Context, name string, values []interface{}) error {
	fmt.Printf("%#v\n", values)
	ctx, collection, cancel, err := connectCollection(ctx, name)
	if err != nil {
		return err
	}
	defer cancel()
	err = collection.Drop(ctx)
	if err != nil {
		return err
	}
	_, err = collection.InsertMany(ctx, values)
	return err
}
