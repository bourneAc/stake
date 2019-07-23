#!/bin/sh

set -e

if [ ! -f "build/env.sh" ]; then
    echo "$0 must be run from the root of the repository."
    exit 2
fi

# Create fake Go workspace if it doesn't exist yet.
workspace="$PWD/build/_workspace"
root="$PWD"
serodir="$workspace/src/github.com/sero-cash"
if [ ! -L "$serodir/stake" ]; then
    mkdir -p "$serodir"
    cd "$serodir"
    ln -s ../../../../../. stake
    cd "$root"
fi

# Set up the environment to use the workspace.
# Also add Godeps workspace so we build using canned dependencies.
GOPATH="$workspace"
GOBIN="$PWD/build/bin"
export GOPATH GOBIN

# Run the command inside the workspace.
cd "$serodir/stake"
PWD="$serodir/stake"

# Launch the arguments with the configured environment.
exec "$@"