package main

import (
	"bufio"
	"log"
	"os"
)

func createRoadmap(inputFile string) (Project, error) {
	lines, err := readRoadmap(inputFile)
	if err != nil {
		return Project{}, err
	}

	roadmap, err := parseRoadmap(lines)
	if err != nil {
		return Project{}, err
	}

	r := roadmap.ToPublic(roadmap.GetFrom(), roadmap.GetTo())

	return r, nil
}

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
