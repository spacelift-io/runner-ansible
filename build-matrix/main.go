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

// amazonAWSVersions maps an ansible major to the amazon.aws requirement the aws image installs.
//
// A galaxy install writes into ANSIBLE_COLLECTIONS_PATH, which takes precedence over the copy
// bundled in the ansible package, and ansible-galaxy does not check the ansible-core that a
// collection asks for. A wrong pair therefore reaches users instead of failing the build.
//
// Each ansible major maps to the amazon.aws major that suits its ansible-core, as an open range
// inside that major. Below 12 the bundled collection is older than the range, so the install runs
// and the collection lands in ANSIBLE_COLLECTIONS_PATH. From 12 on, every minor already bundles a
// collection inside the range, so the install finds nothing to do and the bundled copy stays. Keep
// the ranges open: ansible only moves a bundled collection forward inside its major, and an exact
// version would start shadowing the newer bundled copy as soon as it did.
//
// Ansible 11 is the one entry that installs past what ansible ships. It bundles amazon.aws 9.5.2,
// and 9.5.2 fails on the ansible-core 2.18.19 of the same release: the aws_ec2 inventory plugin
// raises a KeyError on 'version' while it resolves its options, so no inventory is parsed. There is
// no newer 9.x, so the range has to move to 10.
//
// An ansible major that is missing here gets an empty requirement, and the aws build then fails and
// asks for a new entry. Ansible 7 never installs anything, its entry only keeps the map total.
// See scripts/install-collection.sh.
var amazonAWSVersions = map[uint64]string{
	7:  "amazon.aws:>=9.0.0,<10.0.0",
	8:  "amazon.aws:>=9.0.0,<10.0.0",
	9:  "amazon.aws:>=9.0.0,<10.0.0",
	10: "amazon.aws:>=9.0.0,<10.0.0",
	11: "amazon.aws:>=10.0.0,<11.0.0",
	12: "amazon.aws:>=10.0.0,<11.0.0",
	13: "amazon.aws:>=10.0.0,<11.0.0",
	14: "amazon.aws:>=11.0.0,<12.0.0",
}

type ReleaseResponse struct {
	Releases map[string]any `json:"releases"`
}

type matrixVersion struct {
	Ansible        string   `json:"ansible"`
	AdditionalTags []string `json:"additional_tags"`
	AmazonAWS      string   `json:"amazon_aws"`
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
				AmazonAWS:      amazonAWSVersions[version.Major()],
			})
		}
	}

	return matrix
}
