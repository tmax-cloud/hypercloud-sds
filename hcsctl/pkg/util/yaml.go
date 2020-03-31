/*
Usage:
GetValueFromYamlFile(YamlFILE, KIND, YamlKeyPATH)
values, err := GetValueFromYamlFile("/path/to/cluster.yaml", "CephCluster", "metadata.name")
if err != nil {
	panic(err)
}

for _, val := range values{
	fmt.Println(val)
	// convert to string
	val.(string)
	// OR
	str := fmt.Sprintf("%v", val)
}

*/

package util

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

// GetYamlByKindFromByte get all yaml documents from []byte (yaml file)
func GetYamlByKindFromByte(yamlByte []byte, myKind string) ([]yaml.MapSlice, error) {
	var myYAMLs []yaml.MapSlice

	dec := yaml.NewDecoder(bytes.NewReader(yamlByte))

	var oneYamlDoc yaml.MapSlice
	for dec.Decode(&oneYamlDoc) == nil {
		for _, item := range oneYamlDoc {
			key, ok := item.Key.(string)
			if !ok || key != "kind" {
				continue
			}

			value, ok := item.Value.(string)
			if !ok {
				return nil, errors.New("Cannot convert (value) '" +
					fmt.Sprintf("%v", item.Value) + "' to string")
			}

			if value == myKind {
				myYAMLs = append(myYAMLs, oneYamlDoc)
				break
			}
		}
	}

	if len(myYAMLs) == 0 {
		return nil, errors.New("NOT FOUND '" + myKind + "'")
	}

	return myYAMLs, nil
}

// getYamlItemWithoutIndex get yaml item match myKey
func getYamlItemWithoutIndex(myYAML yaml.MapSlice, myKey string) (yaml.MapItem, error) {
	for _, item := range myYAML {
		key := fmt.Sprintf("%v", item.Key)
		if key == myKey {
			return item, nil
		}
	}
	return yaml.MapItem{Key: nil, Value: nil},
		errors.New("NOT FOUND '" + myKey + "'")
}

// getYamlItemWithArrayIndex get yaml item match key with index (ex: names[2])
func getYamlItemWithArrayIndex(myYAML yaml.MapSlice, myKeyIndex string) (yaml.MapItem, error) {
	// Get only NAME from 'NAME[?]'
	myKeyStr := strings.Split(myKeyIndex, "[")[0]

	// Get only NUMBER from '?[NUMBER]'
	myIndexStr := strings.TrimSuffix(strings.Split(myKeyIndex, "[")[1], "]")

	// Convert Index string to int
	myIndexInt, err := strconv.Atoi(myIndexStr)
	if err != nil {
		return yaml.MapItem{Key: nil, Value: nil}, err
	}

	for _, item := range myYAML {
		itemKey := fmt.Sprintf("%v", item.Key)

		if itemKey == myKeyStr {
			myValue, ok := item.Value.([]interface{})
			if !ok {
				return yaml.MapItem{
						Key:   myKeyIndex,
						Value: nil,
					}, errors.New("Type of '" +
						fmt.Sprintf("%v", item.Value) +
						"' is not array []interface{}")
			}

			if len(myValue) == 0 {
				return yaml.MapItem{
					Key:   myKeyIndex,
					Value: nil,
				}, nil
			}

			for index, val := range myValue {
				if index == myIndexInt {
					return yaml.MapItem{
						Key:   myKeyIndex,
						Value: val,
					}, nil
				}
			}
		}
	}

	return yaml.MapItem{Key: nil, Value: nil},
		errors.New("NOT FOUND '" + myKeyIndex + "'")
}

// getYamlItem get yaml item match keyPath (ex: metadata.namespace)
func getYamlItem(myYAML yaml.MapSlice, keyPath string) (yaml.MapItem, error) {
	var item yaml.MapItem

	for _, key := range strings.Split(keyPath, ".") {
		match, err := regexp.MatchString(`^(.*?)\[\d*?\]`, key)
		if err != nil {
			return yaml.MapItem{Key: nil, Value: nil}, err
		}

		if match {
			item, err = getYamlItemWithArrayIndex(myYAML, key)
		} else {
			item, err = getYamlItemWithoutIndex(myYAML, key)
		}

		if err != nil {
			return item, err
		}

		myYAML, _ = item.Value.(yaml.MapSlice)
	}

	return item, nil
}

// GetValueFromYamlByte get all values of keyPath from []byte (yaml file)
func GetValueFromYamlByte(yamlByte []byte, kind, keyPath string) ([]interface{}, error) {
	allYAMLs, err := GetYamlByKindFromByte(yamlByte, kind)
	if err != nil {
		return nil, err
	}

	var itemValues []interface{}
	for _, oneYAML := range allYAMLs {
		item, err := getYamlItem(oneYAML, keyPath)
		if err != nil {
			continue
		}
		itemValues = append(itemValues, item.Value)
	}

	return itemValues, nil
}

// GetValueFromYamlFile get all values of keyPath from yaml file
func GetValueFromYamlFile(filename, kind, keyPath string) ([]interface{}, error) {
	fileByte, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return GetValueFromYamlByte(fileByte, kind, keyPath)
}
