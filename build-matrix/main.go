package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"

	"github.com/Masterminds/semver/v3"
)

const (
	// Define the oldest major version we care about, we do not want to build image starting ansible 1.0
	minSupportedMajor = 7
	// Numbed of latest minor release to keep building
	maxSupportedMinor = 5
)

type ReleaseResponse struct {
	Releases map[string]any `json:"releases"`
}

type matrixVersion struct {
	Ansible        string   `json:"ansible"`
	AdditionalTags []string `json:"additional_tags"`
}
type Matrix []matrixVersion

// This small script reads ansible versions from pypi and returns an aggregated list of deduplicated minor version.
// This is used to compute the build matrix in github to build all minor versions in parallel.
// This script will also find the latest minor version for every major to be able to tag the docker image dynamically.
// For example if we have the following version returned:
// - 10.1.1
// - 10.3.1
// - 10.4.3
// - 10.4.3.
// The script will return the following versions:
// - 10.1
// - 10.3
// - 10.4, additional_tags: 10
// We also only keep the latest 5 minor versions because we are limited to 256 jobs per workflow run
// Check main_test.go for a quick overview of the expected behavior.
func main() {
	resp, err := http.Get("https://pypi.org/pypi/ansible/json")
	if err != nil {
		log.Fatal(err)
	}

	matrixOutput := GenerateBuildMatrix(resp.Body, minSupportedMajor)
	output, err := json.Marshal(matrixOutput)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(string(output))
}

func GenerateBuildMatrix(reader io.Reader, minSupportedMajor uint64) Matrix {
	releases := ReleaseResponse{}
	if err := json.NewDecoder(reader).Decode(&releases); err != nil {
		log.Fatal(err)
	}

	var versions []*semver.Version

	for v := range releases.Releases {
		version, err := semver.NewVersion(v)
		if err != nil {
			log.Printf("Unable to parse version %s\n", v)
			continue
		}
		versions = append(versions, version)
	}

	sort.Slice(versions, func(i, j int) bool {
		return versions[j].LessThan(versions[i])
	})

	versionGroupedByMajor := make(map[int][]*semver.Version)
	// Just used for stable ordering
	var majorVersions []int

	for _, version := range versions {
		if version.Major() < minSupportedMajor {
			break
		}
		major := int(version.Major())
		if _, exists := versionGroupedByMajor[major]; !exists {
			majorVersions = append(majorVersions, major)
		}
		if len(versionGroupedByMajor[major]) < maxSupportedMinor {
			versionGroupedByMajor[major] = append(versionGroupedByMajor[major], version)
		}
	}

	sort.Sort(sort.Reverse(sort.IntSlice(majorVersions)))

	minorVersionDeduplication := map[string]any{}
	matrix := Matrix{}
	for j, majorVersion := range majorVersions {
		for i, version := range versionGroupedByMajor[majorVersion] {
			key := fmt.Sprintf("%d.%d", version.Major(), version.Minor())
			if _, exists := minorVersionDeduplication[key]; exists {
				continue
			}
			additionalTags := make([]string, 0)
			if i == 0 {
				additionalTags = append(additionalTags, fmt.Sprintf("%d", version.Major()))
				if j == 0 {
					additionalTags = append(additionalTags, "latest")
				}
			}

			minorVersionDeduplication[key] = struct{}{}
			matrix = append(matrix, matrixVersion{
				Ansible:        key,
				AdditionalTags: additionalTags,
			})
		}
	}

	return matrix
}
