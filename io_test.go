package main

import (
	"reflect"
	"testing"
)

func Test_createRoadmap(t *testing.T) {
	type args struct {
		inputFile string
	}
	tests := []struct {
		name    string
		args    args
		want    Project
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := createRoadmap(tt.args.inputFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("createRoadmap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createRoadmap() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_readRoadmap(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readRoadmap(tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("readRoadmap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readRoadmap() got = %v, want %v", got, tt.want)
			}
		})
	}
}
