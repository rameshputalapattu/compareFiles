package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"log"
	"os"
	"path/filepath"
)

func main() {

	firstDir := os.Args[1]
	secondDir := os.Args[2]

	fileset1, err := getFilesInDir(firstDir)
	if err != nil {
		log.Fatal(err)
	}

	fileset2, err := getFilesInDir(secondDir)
	if err != nil {
		log.Fatal(err)
	}

	h := sha256.New()
	digestMap1, err := getContentDigestMap(firstDir, fileset1, h)
	if err != nil {
		log.Fatal(err)
	}

	digestMap2, err := getContentDigestMap(secondDir, fileset2, h)

	if err != nil {
		log.Fatal(err)
	}

	compareDigestMaps(digestMap1, digestMap2)

}

func getContentDigestMap(dirName string, fileNames []string, h hash.Hash) (map[string]string, error) {
	digestMap := make(map[string]string)

	for _, fname := range fileNames {
		fullName := filepath.Join(dirName, fname)
		f, err := os.Open(fullName)
		if err != nil {
			return nil, err
		}

		h.Reset()
		_, err = io.Copy(h, f)

		if err != nil {
			return nil, err
		}

		digestMap[fname] = hex.EncodeToString(h.Sum(nil))

	}
	return digestMap, nil
}

func getFilesInDir(dirName string) ([]string, error) {
	fdir, err := os.Open(dirName)
	defer fdir.Close()
	if err != nil {

		return nil, err

	}

	var filesInDir []string
	fileinfo, err := fdir.Readdir(-1)
	if err != nil {
		return nil, err
	}

	for _, finfo := range fileinfo {
		if !finfo.IsDir() {
			filesInDir = append(filesInDir, finfo.Name())
		}
	}

	return filesInDir, nil

}

func compareDigestMaps(digestMap1 map[string]string, digestMap2 map[string]string) {
	fmt.Println("=================================================")
	fmt.Println("fileName|set1|set2|matched?")
	fmt.Println("=================================================")
	var allfileNames []string
	for key := range digestMap1 {
		allfileNames = append(allfileNames, key)
	}

	for key := range digestMap2 {
		_, ok := digestMap1[key]
		if !ok {
			allfileNames = append(allfileNames, key)
		}
	}

	for _, fname := range allfileNames {
		digest1 := digestMap1[fname]
		digest2 := digestMap2[fname]
		matched := digest1 == digest2
		fmt.Printf("%s|%s|%s|%t\n", fname, digest1, digest2, matched)
	}
}
