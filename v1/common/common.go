package common

const (
	ROOT = ".yogit"

	ROOT_DIR_OBJECTS = ".yogit/objects"
	ROOT_DIR_LOG     = ".yogit/log"
	ROOT_DIR_REFS    = ".yogit/refs"

	ROOT_STAGE_FILE   = ".yogit/stage"
	ROOT_GLOBALS_FILE = ".yogit/GLOBALS"
	ROOT_HEADER_FILE  = ".yogit/HEADER"

	BRANCH_PATH   = ".yogit/refs/heads"
	MASTER_BRANCH = ".yogit/refs/heads/master"
	BRANCH_REFS   = "refs:refs/heads"
	LOG_LOGS_FILE = ".yogit/log/logs.txt"
)

var IGNORE_FILES = []string{
	".git", ".png", ".jpg", "__pycache__",
	"node_modules", ".yogit",
}

// Represents the hashId for commits
type CommitId string
