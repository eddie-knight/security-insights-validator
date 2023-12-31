package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

// Validate validates the SecurityInsightsSchema by
// 1. Unmarshalling the input file into the SecurityInsightsSchema
// 2. Marshalling the SecurityInsightsSchema into a new file
// 3. Diffing the input file and the new file
func (s *SecurityInsightsSchema) Validate() (err error) {
	err = s.ingestFile()
	if err != nil {
		return
	}
	err = s.compareData()
	if err != nil {
		return
	}
	return nil
}

// IngestFile reads the input file and unmarshals it into the SecurityInsightsSchema
func (s *SecurityInsightsSchema) ingestFile() error {
	yamlFile, err := ioutil.ReadFile(viper.GetString("input"))
	if err != nil {
		return reformatError("Error Reading Provided File", err)
	}
	err = yaml.Unmarshal(yamlFile, s)
	if err != nil {
		return reformatError("Error Unmarshalling to Specification", err)
	}
	return nil
}

// compareData compares the input file to the marshalled SecurityInsightsSchema
func (s *SecurityInsightsSchema) compareData() error {
	// marshal the SecurityInsightsSchema into bytes
	yamlBytes, err := yaml.Marshal(s)
	if err != nil {
		return err
	}

	// open the input file as bytes
	inputFile, err := ioutil.ReadFile(viper.GetString("input"))
	if err != nil {
		return err
	}

	// Use the diff function to compare the input file to the unmarshalled SecurityInsightsSchema
	diff, err := diff("clean", yamlBytes, inputFile)
	fmt.Print(string(diff))
	if err != nil {
		return err
	}
	return nil
}

func init() {
	NewFileName = fmt.Sprintf("clean-%s", viper.GetString("input"))
}

// Returns diff of two arrays of bytes in diff tool format.
// ref: https://cs.opensource.google/go/go/+/refs/tags/go1.18.10:src/cmd/internal/diff/diff.go
func diff(prefix string, b1, b2 []byte) ([]byte, error) {
	writeTempFile := func(prefix string, data []byte) (string, error) {
		file, err := ioutil.TempFile("", prefix)
		if err != nil {
			return "", err
		}
		_, err = file.Write(data)
		if err1 := file.Close(); err == nil {
			err = err1
		}
		if err != nil {
			os.Remove(file.Name())
			return "", err
		}
		return file.Name(), nil
	}

	f1, err := writeTempFile(prefix, b1)
	if err != nil {
		return nil, err
	}
	defer os.Remove(f1)

	f2, err := writeTempFile(prefix, b2)
	if err != nil {
		return nil, err
	}
	defer os.Remove(f2)

	cmd := "diff"
	if runtime.GOOS == "plan9" {
		cmd = "/bin/ape/diff"
	}

	data, err := exec.Command(cmd, "-u", f1, f2).CombinedOutput()
	if len(data) > 0 {
		// diff exits with a non-zero status when the files don't match.
		// Ignore that failure as long as we get output.
		err = nil
	}

	// If we are on Windows and the diff is Cygwin diff,
	// machines can get into a state where every Cygwin
	// command works fine but prints a useless message like:
	//
	//	Cygwin WARNING:
	//	  Couldn't compute FAST_CWD pointer.  This typically occurs if you're using
	//	  an older Cygwin version on a newer Windows.  Please update to the latest
	//	  available Cygwin version from https://cygwin.com/.  If the problem persists,
	//	  please see https://cygwin.com/problems.html
	//
	// Skip over that message and just return the actual diff.
	if len(data) > 0 && !bytes.HasPrefix(data, []byte("--- ")) {
		i := bytes.Index(data, []byte("\n--- "))
		if i >= 0 && i < 80*10 && bytes.Contains(data[:i], []byte("://cygwin.com/")) {
			data = data[i+1:]
		}
	}

	return data, err
}
