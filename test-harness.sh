#!/usr/bin/env bash
set -euo pipefail

# test-harness.sh setup/start cucumber test environment.
#
# Configuration is managed with environment variables, the ones you
# are most likely to reconfigured are stored in '.test-env'.
#
# Variables:
#   SDK_TESTING_URL     - URL to algorand-sdk-testing, useful for forks.
#   SDK_TESTING_BRANCH  - branch to checkout, useful for new tests.
#   SDK_TESTING_HARNESS - local directory that the algorand-sdk-testing repo is cloned into.
#   VERBOSE_HARNESS     - more output while the script runs.
#   INSTALL_ONLY        - installs feature files only, useful for unit tests.
#
#   WARNING: If set to 1, new features will be LOST when downloading the test harness.
#   REGARDLESS: modified features are ALWAYS overwritten.
#   REMOVE_LOCAL_FEATURES - delete all local cucumber feature files before downloading these from github.
#
#   WARNING: Be careful when turning on the next variable.
#   In that case you'll need to provide all variables expected by `algorand-sdk-testing`'s `.env`
#   OVERWRITE_TESTING_ENVIRONMENT=0
#
#   LOCAL_SDK_BUILD    - build the harness's indexer and conduit against this working tree.
#   STAGED_SANDBOX     - local directory that the patched sandbox is staged in.

SHUTDOWN=0
if [ $# -ne 0 ]; then
  if [ $# -ne 1 ]; then
    echo "this script accepts a single argument, which must be 'up' or 'down'."
    exit 1
  fi

  case $1 in
    'up')
      ;; # default.
    'down')
      SHUTDOWN=1
      ;;
    *)
      echo "unknown parameter '$1'."
      echo "this script accepts a single argument, which must be 'up' or 'down'."
      exit 1
      ;;
  esac
fi

START=$(date "+%s")

THIS=$(basename "$0")
ENV_FILE=".test-env"
TEST_DIR="test"

set -a
source "$ENV_FILE"
set +a

rootdir=$(dirname "$0")
pushd "$rootdir"

echo "$THIS: VERBOSE_HARNESS=$VERBOSE_HARNESS"

# Read one variable from algorand-sdk-testing's own .env, without letting its
# settings leak into ours (it defines VERBOSE_HARNESS too, for instance).
harness_env() {
  ( set -a; source "$SDK_TESTING_HARNESS"/.env; set +a; printf '%s' "${!1}" )
}

GO_MOD_REPLACE='go mod edit -replace github.com/algorand/go-algorand-sdk/v2=/tmp/local-sdk'

# Rewrite the line matching $2 as the perl replacement $3, then confirm $4 is in
# the file. These patches land in files owned by the sandbox repo, so a silent
# no-op is a real possibility -- and would resurface much later as a baffling
# failure from deep inside a docker build.
patch_line() {
  local file=$1 pattern=$2 replacement=$3 expect=$4
  perl -0pi -e "s{^$pattern\$}{$replacement}m" "$file"
  if ! grep -qF "$expect" "$file"; then
    echo "$THIS: cannot patch $file, no line matches /^$pattern\$/" >&2
    exit 1
  fi
}

# Give one sandbox image the staged SDK and build it against that instead of the
# release its go.mod pins. A directory replace covers the modules that reach the
# SDK indirectly too, so patching conduit is enough to also rebuild the indexer
# packages it embeds. $3 overrides the build command that install.sh runs.
patch_sandbox_image() {
  local dir=$1 dockerfile=$2 build=${3:-make}
  patch_line "$dir/$dockerfile" 'RUN /tmp/install\.sh' \
    'COPY images/local-sdk /tmp/local-sdk\n$&' 'COPY images/local-sdk'
  # Conduit has its own go.sum lacking any SDK changes, let's tidy it up
  patch_line "$dir/install.sh" 'make' "$GO_MOD_REPLACE\ngo mod tidy\n$build" "$GO_MOD_REPLACE"
}

# The sandbox builds indexer and conduit by cloning them inside a container, so
# each compiles against the go-algorand-sdk release its go.mod pins. Whenever
# this repo adds support for something those releases predate -- a new consensus
# version, most often -- the pinned builds reject the blocks algod produces and
# the harness never comes up. Staging our own copy of the sandbox lets both
# builds compile against this working tree, with no need to publish an SDK
# branch for them to depend on.
stage_local_sdk_sandbox() {
  local url branch sandbox
  url=$(harness_env SANDBOX_URL)
  branch=$(harness_env SANDBOX_BRANCH)
  sandbox="$(pwd)/$STAGED_SANDBOX"

  git clone --depth 1 --single-branch --branch "$branch" "$url" "$sandbox"

  # The sandbox's .dockerignore excludes everything outside a few directories,
  # so the SDK has to be staged under images/ to reach the build context.
  rsync -a --exclude .git --exclude "$SDK_TESTING_HARNESS" --exclude "$STAGED_SANDBOX" \
    ./ "$sandbox"/images/local-sdk

  patch_sandbox_image "$sandbox"/images/indexer IndexerDockerfile

  # Conduit's default make target regenerates its filter-processor field map by
  # reflecting over the SDK's transaction types, and that generator treats a
  # codec tag verbatim when it recurses, so an option like `txn,required` throws
  # off every nested tag path it looks up. It also has no cast for the pqsig
  # fields. Both leave it rejecting any SDK newer than the release conduit pins.
  # The file it would rewrite is committed and still compiles, so build the
  # binary and leave generation out of it.
  patch_sandbox_image "$sandbox"/images/conduit Dockerfile '(cd cmd/conduit && go build)'

  # up.sh gets at the sandbox by cloning it, which only carries committed state.
  git -C "$sandbox" add -Af
  git -C "$sandbox" \
    -c user.name="$THIS" -c user.email="$THIS@localhost" \
    commit -qm "build indexer and conduit against the local go-algorand-sdk"

  # up.sh sources .env with `set -a`, so appending overrides the defaults that
  # algorand-sdk-testing ships with.
  {
    echo "SANDBOX_URL=\"file://$sandbox\""
    echo "SANDBOX_BRANCH=\"$branch\""
  } >> "$SDK_TESTING_HARNESS"/.env
}

rm -rf "$STAGED_SANDBOX"

## Reset test harness
if [ -d "$SDK_TESTING_HARNESS" ]; then
  pushd "$SDK_TESTING_HARNESS"
  ./scripts/down.sh
  popd
  rm -rf "$SDK_TESTING_HARNESS"
  if [[ $SHUTDOWN == 1 ]]; then
    echo "$THIS: network shutdown complete."
    exit 0
  fi
else
  echo "$THIS: directory $SDK_TESTING_HARNESS does not exist - NOOP"
fi

if [[ $SHUTDOWN == 1 ]]; then
  echo "$THIS: unable to shutdown network."
  exit 1
fi

git clone --depth 1 --single-branch --branch "$SDK_TESTING_BRANCH" "$SDK_TESTING_URL" "$SDK_TESTING_HARNESS"

echo "$THIS: OVERWRITE_TESTING_ENVIRONMENT=$OVERWRITE_TESTING_ENVIRONMENT"
if [[ $OVERWRITE_TESTING_ENVIRONMENT == 1 ]]; then
  echo "$THIS: OVERWRITE replaced $SDK_TESTING_HARNESS/.env with $ENV_FILE:"
  cp "$ENV_FILE" "$SDK_TESTING_HARNESS"/.env
fi

echo "$THIS: REMOVE_LOCAL_FEATURES=$REMOVE_LOCAL_FEATURES"
## Copy feature files into the project resources
if [[ $REMOVE_LOCAL_FEATURES == 1 ]]; then
  echo "$THIS: OVERWRITE wipes clean $TEST_DIR/features"
  if [[ $VERBOSE_HARNESS == 1 ]]; then
    ( tree $TEST_DIR/features && echo "$THIS: see the previous for files deleted" ) || true
  fi
  rm -rf $TEST_DIR/features
fi
mkdir -p $TEST_DIR/features
cp -r "$SDK_TESTING_HARNESS"/features/* $TEST_DIR/features
if [[ $VERBOSE_HARNESS == 1 ]]; then
  ( tree $TEST_DIR/features && echo "$THIS: see the previous for files copied over" ) || true
fi
echo "$THIS: seconds it took to get to end of cloning and copying: $(($(date "+%s") - START))s"

if [[ $INSTALL_ONLY == 1 ]]; then
  echo "$THIS: configured to install feature files only. Not starting test harness environment."
  exit 0
fi

echo "$THIS: LOCAL_SDK_BUILD=$LOCAL_SDK_BUILD"
if [[ $LOCAL_SDK_BUILD == 1 ]]; then
  stage_local_sdk_sandbox
  echo "$THIS: seconds it took to stage $STAGED_SANDBOX: $(($(date "+%s") - START))s"
fi

## Start test harness environment
pushd "$SDK_TESTING_HARNESS"

[[ "$VERBOSE_HARNESS" = 1 ]] && V_FLAG="-v" || V_FLAG=""
echo "$THIS: standing up harnness with command [./up.sh $V_FLAG]"
./scripts/up.sh "$V_FLAG"

popd
echo "$THIS: seconds it took to finish testing sdk's up.sh: $(($(date "+%s") - START))s"
echo ""
echo "--------------------------------------------------------------------------------"
echo "|"
echo "|    To run sandbox commands, cd into $SDK_TESTING_HARNESS/.sandbox             "
echo "|"
echo "--------------------------------------------------------------------------------"
