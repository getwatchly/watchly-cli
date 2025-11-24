TEMP_DIR=$(mktemp -d)
cd $TEMP_DIR

semver_regex='^(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)(-(0|[1-9][0-9]*|[0-9]*[a-zA-Z-][0-9a-zA-Z-]*)(\.(0|[1-9][0-9]*|[0-9]*[a-zA-Z-][0-9a-zA-Z-]*))*)?(\+([0-9a-zA-Z-]+(\.[0-9a-zA-Z-]+)*))?$'

validate_semver() {
    local version_string="$1"
    if [[ "$version_string" =~ $semver_regex ]]; then
        return 0
    else
        return 1
    fi
}

CLI_VERSION=$WATCHLY_CLI_VERSION
if [ -z "$CLI_VERSION" ]; then
    CLI_VERSION="0.0.9"
fi

curl -L https://github.com/getwatchly/watchly-cli/releases/download/v$CLI_VERSION/watchly-cli_${CLI_VERSION}_linux_386.tar.gz > watchly-cli.tar.gz
tar -xzf watchly-cli.tar.gz
cp watchly-cli /usr/local/bin

cd $HOME
rm -r $TEMP_DIR

echo "Watchly CLI version $CLI_VERSION installed successfully!"
