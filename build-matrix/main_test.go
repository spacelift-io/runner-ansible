package main

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateBuildMatrix(t *testing.T) {
	fakePythonVersions := ReleaseResponse{
		Releases: map[string]any{
			"1.1.0":  struct{}{},
			"1.1.1":  struct{}{},
			"2.10.0": struct{}{},
			"2.11.0": struct{}{},
			"2.11.2": struct{}{},
			"3.1.0":  struct{}{},
			"3.1.1":  struct{}{},
			"3.2.0":  struct{}{},
			"3.3.0":  struct{}{},
			"3.4.0":  struct{}{},
			"3.5.0":  struct{}{},
			"3.6.0":  struct{}{},
		},
	}

	fakeJsonResponse, err := json.Marshal(fakePythonVersions)
	require.NoError(t, err)

	matrix := GenerateBuildMatrix(bytes.NewReader(fakeJsonResponse), 2)
	expectedMatrix := Matrix{
		{
			Ansible:        "3.6",
			AdditionalTags: []string{"3"},
		},
		{
			Ansible:        "3.5",
			AdditionalTags: []string{},
		},
		{
			Ansible:        "3.4",
			AdditionalTags: []string{},
		},
		{
			Ansible:        "3.3",
			AdditionalTags: []string{},
		},
		{
			Ansible:        "3.2",
			AdditionalTags: []string{},
		},
		{
			Ansible:        "2.11",
			AdditionalTags: []string{"2"},
		},
		{
			Ansible:        "2.10",
			AdditionalTags: []string{},
		},
	}
	assert.Equal(t, expectedMatrix, matrix)
}
