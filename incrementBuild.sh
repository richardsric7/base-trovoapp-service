#!/bin/bash

set -e  # Exit on errors
set -x  # Print commands as they are executed

path_to_pubspec="pubspec.yaml"
current_version=$(awk '/^version:/ {print $2}' $path_to_pubspec)
echo "current version: $current_version"
current_build_version=$(echo "$current_version" | sed 's/.*+//')
echo "current build version: $current_build_version"
current_version_without_build=$(echo "$current_version" | sed 's/\+.*//')
echo "current version without build: $current_version_without_build"
trimmed_string="${current_build_version#"${current_build_version%%[![:space:]]*}"}"
trimmed_string="${trimmed_string%"${trimmed_string##*[![:space:]]}"}"
build_increment=$(expr $trimmed_string + 1)
echo "build increment: $build_increment"
new_version="$current_version_without_build+$build_increment"
echo "new version: $new_version"
echo "Setting pubspec.yaml version $current_version to $new_version"
if [[ "$OSTYPE" == "darwin"* ]]; then
    # macOS sed (requires a space after -i)
    sed -i '' -e "s/version: $current_version/version: $new_version/g" $path_to_pubspec
else
    # GNU sed (requires no space after -i)
    sed -i'' -e "s/version: $current_version/version: $new_version/g" $path_to_pubspec
fi
