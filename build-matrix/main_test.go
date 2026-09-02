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
			AdditionalTags: []string{"3", "latest"},
			AmazonAWS:      "",
		},
		{
			Ansible:        "3.5",
			AdditionalTags: []string{},
			AmazonAWS:      "",
		},
		{
			Ansible:        "3.4",
			AdditionalTags: []string{},
			AmazonAWS:      "",
		},
		{
			Ansible:        "3.3",
			AdditionalTags: []string{},
			AmazonAWS:      "",
		},
		{
			Ansible:        "3.2",
			AdditionalTags: []string{},
			AmazonAWS:      "",
		},
		{
			Ansible:        "2.11",
			AdditionalTags: []string{"2"},
			AmazonAWS:      "",
		},
		{
			Ansible:        "2.10",
			AdditionalTags: []string{},
			AmazonAWS:      "",
		},
	}
	assert.Equal(t, expectedMatrix, matrix)
}

func TestGenerateBuildMatrix_AmazonAWS(t *testing.T) {
	fakePythonVersions := ReleaseResponse{
		Releases: map[string]any{
			"10.7.0": struct{}{},
			"13.8.0": struct{}{},
			"14.3.1": struct{}{},
			// A major nobody has mapped yet. It must come out empty, which fails the aws build.
			"15.0.0": struct{}{},
		},
	}

	fakeJsonResponse, err := json.Marshal(fakePythonVersions)
	require.NoError(t, err)

	matrix := GenerateBuildMatrix(bytes.NewReader(fakeJsonResponse), 10)

	expectedRequirements := map[string]string{
		"10.7": amazonAWSVersions[10],
		"13.8": amazonAWSVersions[13],
		"14.3": amazonAWSVersions[14],
		"15.0": "",
	}

	require.Len(t, matrix, len(expectedRequirements))
	for _, version := range matrix {
		expected, found := expectedRequirements[version.Ansible]
		require.True(t, found, "unexpected ansible version %s in the matrix", version.Ansible)
		assert.Equal(t, expected, version.AmazonAWS, "amazon.aws requirement of ansible %s", version.Ansible)
	}
}

// TestAmazonAWSVersionsCoverSupportedMajors fails when someone raises minSupportedMajor or a new
// ansible major appears without an entry in the map.
func TestAmazonAWSVersionsCoverSupportedMajors(t *testing.T) {
	for major := uint64(minSupportedMajor); major <= 14; major++ {
		assert.NotEmpty(t, amazonAWSVersions[major], "ansible %d has no amazon.aws requirement", major)
	}
}
