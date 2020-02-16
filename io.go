package main

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"
)

func readRoadmap(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func writeRoadmap(filename, content string) error {
	d1 := []byte(content)
	err := ioutil.WriteFile(filename, d1, 0644)

	return err
}
