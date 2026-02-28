all: test

MORES_FILE_PREFIX=../../storage/
MORES_FILE_SUFFIX=.txt
MORES_PORT=8002

export MORES_FILE_PREFIX MORES_FILE_SUFFIX MORES_PORT

# Extract the optional argument after "test" from the command line
EXTRA_ARGS := $(filter-out test,$(MAKECMDGOALS))
ifeq ($(EXTRA_ARGS),)
TEST_SUBDIR :=
else
TEST_SUBDIR := $(firstword $(EXTRA_ARGS))
endif

# Treat "all" as no specific subdirectory (i.e., run all tests)
ifeq ($(TEST_SUBDIR),all)
TEST_SUBDIR :=
endif

# Create dummy targets for any extra arguments to avoid "No rule to make target" errors
$(EXTRA_ARGS):
	@# do nothing

test:
	if [ -z "$(TEST_SUBDIR)" ]; then \
		go test -v ./tests/...; \
	else \
		go test -v ./tests/$(TEST_SUBDIR)/...; \
	fi
