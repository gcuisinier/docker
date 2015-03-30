package main

import (
	"io/ioutil"
	"os/exec"
	"strings"
	"testing"
)

func TestImportDisplay(t *testing.T) {
	runCmd := exec.Command(dockerBinary, "run", "-d", "busybox", "true")
	out, _, err := runCommandWithOutput(runCmd)
	if err != nil {
		t.Fatal("failed to create a container", out, err)
	}
	cleanedContainerID := stripTrailingCharacters(out)
	defer deleteContainer(cleanedContainerID)

	out, _, err = runCommandPipelineWithOutput(
		exec.Command(dockerBinary, "export", cleanedContainerID),
		exec.Command(dockerBinary, "import", "-"),
	)
	if err != nil {
		t.Errorf("import failed with errors: %v, output: %q", err, out)
	}

	if n := strings.Count(out, "\n"); n != 1 {
		t.Fatalf("display is messed up: %d '\\n' instead of 1:\n%s", n, out)
	}
	image := strings.TrimSpace(out)
	defer deleteImages(image)

	runCmd = exec.Command(dockerBinary, "run", "--rm", image, "true")
	out, _, err = runCommandWithOutput(runCmd)
	if err != nil {
		t.Fatal("failed to create a container", out, err)
	}

	if out != "" {
		t.Fatalf("command output should've been nothing, was %q", out)
	}

	logDone("import - display is fine, imported image runs")
}

func TestImportFile(t *testing.T) {
	runCmd := exec.Command(dockerBinary, "run", "-d", "busybox", "true")
	out, _, err := runCommandWithOutput(runCmd)
	if err != nil {
		t.Fatal("failed to create a container", out, err)
	}
	cleanedContainerID := stripTrailingCharacters(out)
	defer deleteContainer(cleanedContainerID)

	runCmd = exec.Command(dockerBinary, "export", cleanedContainerID)
	out, _, err = runCommandWithOutput(runCmd)
	if err != nil {
		t.Fatal("failed to export a container", out, err)
	}

	temporaryFile, err := ioutil.TempFile("", "exportImportTest")
	if err != nil {
		t.Fatal("failed to create temporary file", "", err)
	}
	temporaryFile.WriteString(out)

	runCmd = exec.Command(dockerBinary, "import", temporaryFile.Name())
	out, _, err = runCommandWithOutput(runCmd)
	if err != nil {
		t.Fatal("failed to import a container", out, err)
	}

	if n := strings.Count(out, "\n"); n != 1 {
		t.Fatalf("display is messed up: %d '\\n' instead of 1:\n%s", n, out)
	}
	image := strings.TrimSpace(out)
	defer deleteImages(image)

	runCmd = exec.Command(dockerBinary, "run", "--rm", image, "true")
	out, _, err = runCommandWithOutput(runCmd)
	if err != nil {
		t.Fatal("failed to create a container", out, err)
	}

	if out != "" {
		t.Fatalf("command output should've been nothing, was %q", out)
	}

	logDone("import file is fine, imported image runs")
}
