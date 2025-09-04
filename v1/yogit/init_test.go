package yogit

import (
	"fmt"
	"os"
	"testing"
	"v1/common"
)

func TestFolderStructure(t *testing.T) {
	tmp := t.TempDir()
	Init(tmp)
	expectedFolders := []string{
		common.ROOT, common.ROOT_DIR_REFS, common.ROOT_DIR_LOG,
		common.ROOT_DIR_OBJECTS, common.ROOT_STAGE_FILE, common.ROOT_HEADER_FILE,
		common.MASTER_BRANCH, common.LOG_LOGS_FILE,
	}

	for _, path := range expectedFolders {
		_, err := os.Stat(path)
		if err != nil {
			t.Errorf("Expected Path of %s, got:\n%v\n", path, err)
		}
	}
}

func TestMasterBranch(t *testing.T) {
	expectedBranch := fmt.Sprintf("%s/%s", common.BRANCH_REFS, "master")
	content, err := os.ReadFile(common.ROOT_HEADER_FILE)
	if err != nil {
		t.Error(err)
	}

	actualBranch := string(content)

	if actualBranch != expectedBranch {
		t.Errorf("The expected branch %s did not meet with %s\n", expectedBranch, actualBranch)
	}
}
