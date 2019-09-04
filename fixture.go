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
	for cn, cd := range ds {
		vs, err := toValues(cn, cd)
		if err != nil {
			return err
		}
		err = resetCollection(ctx, cn, vs)
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
	return "", fmt.Errorf("DataSet %q not found in %s", name, conf.fixtureRootDirAbs)
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

func loadDataSet(files ...string) (DataSet, error) {
	dss, err := toDataSets(files...)
	if err != nil {
		return nil, err
	}
	return mergeDataSet(dss), nil
}

func toDataSets(files ...string) ([]DataSet, error) {
	dss := make([]DataSet, len(files))
	for i, file := range files {
		ds, err := readFixtureFile(file)
		if err != nil {
			return nil, err
		}
		dss[i] = ds
	}
	return dss, nil
}

func readFixtureFile(file string) (DataSet, error) {
	format, err := fixtureFormat(file)
	if err != nil {
		return nil, err
	}
	bs, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	ds := make(DataSet)
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
	if conf.FixtureFormat != FixtureFormatAuto {
		return conf.FixtureFormat, nil
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

func mergeDataSet(dss []DataSet) DataSet {
	merged := make(DataSet)
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

func mergeCollData(coll1, coll2 CollectionData) CollectionData {
	merged := make(CollectionData)
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

func mergeDocData(doc1, doc2 DocData) DocData {
	merged := make(DocData)
	for k, v := range doc1 {
		merged[k] = v
	}
	for k, v := range doc2 {
		merged[k] = v
	}
	return merged
}

func toValues(collectionName string, coll CollectionData) ([]interface{}, error) {
	values := make([]interface{}, 0, len(coll))
	for id, doc := range coll {
		newDoc := make(DocData)
		for k, v := range doc {
			newDoc[k] = v
		}
		newDoc["_id"] = id
		v, err := applyPreFuncs(collectionName, newDoc)
		if err != nil {
			return nil, err
		}
		values = append(values, v)
	}
	return values, nil
}

func applyPreFuncs(collName string, value DocData) (DocData, error) {
	if conf.PreInsertFuncs == nil {
		return value, nil
	}
	v := value
	var err error
	for _, fn := range conf.PreInsertFuncs {
		v, err = fn(collName, v)
		if err != nil {
			return nil, err
		}
	}
	return v, nil
}

func resetCollection(ctx context.Context, name string, values []interface{}) error {
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
