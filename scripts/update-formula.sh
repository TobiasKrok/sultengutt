#!/bin/bash

# Script to manually update Homebrew formula
# Usage: ./update-formula.sh v1.0.0

set -e

VERSION=$1
if [ -z "$VERSION" ]; then
    echo "Usage: $0 <version>"
    echo "Example: $0 v1.0.0"
    exit 1
fi

echo "Fetching release information for $VERSION..."

# Get release data from GitHub API
RELEASE_DATA=$(curl -s https://api.github.com/repos/TobiasKrok/sultengutt/releases/tags/${VERSION})

if [ "$RELEASE_DATA" = "Not Found" ] || [ -z "$RELEASE_DATA" ]; then
    echo "Error: Release ${VERSION} not found"
    exit 1
fi

# Extract URLs
DARWIN_AMD64_URL=$(echo "$RELEASE_DATA" | jq -r '.assets[] | select(.name | contains("darwin-amd64")) | .browser_download_url')
DARWIN_ARM64_URL=$(echo "$RELEASE_DATA" | jq -r '.assets[] | select(.name | contains("darwin-arm64")) | .browser_download_url')

if [ -z "$DARWIN_AMD64_URL" ] || [ -z "$DARWIN_ARM64_URL" ]; then
    echo "Error: Could not find Darwin binaries in release"
    exit 1
fi

echo "Downloading binaries to calculate SHA256..."

# Download and calculate SHA256
curl -L "$DARWIN_AMD64_URL" -o /tmp/darwin-amd64.tar.gz
curl -L "$DARWIN_ARM64_URL" -o /tmp/darwin-arm64.tar.gz

SHA256_AMD64=$(sha256sum /tmp/darwin-amd64.tar.gz | cut -d' ' -f1)
SHA256_ARM64=$(sha256sum /tmp/darwin-arm64.tar.gz | cut -d' ' -f1)

# Clean up
rm /tmp/darwin-amd64.tar.gz /tmp/darwin-arm64.tar.gz

# Generate formula
cat > sultengutt.rb << EOF
class Sultengutt < Formula
  desc "Cross-platform desktop reminder for ordering surprise dinners"
  homepage "https://github.com/TobiasKrok/sultengutt"
  version "${VERSION#v}"
  license "MIT"

  on_macos do
    if Hardware::CPU.intel?
      url "${DARWIN_AMD64_URL}"
      sha256 "${SHA256_AMD64}"
    else
      url "${DARWIN_ARM64_URL}"
      sha256 "${SHA256_ARM64}"
    end
  end

  def install
    bin.install "sultengutt"
  end

  def caveats
    <<~EOS
      To set up Sultengutt, run:
        sultengutt install
      
      To check status:
        sultengutt status
      
      ⚠️  IMPORTANT: Before uninstalling with Homebrew:
        Run 'sultengutt uninstall' first to remove scheduled tasks and config files.
        Then run 'brew uninstall sultengutt' to remove the application.
    EOS
  end

  test do
    assert_match "Sultengutt", shell_output("#{bin}/sultengutt --help")
  end
end
EOF

echo "Formula generated: sultengutt.rb"
echo ""
echo "URLs:"
echo "  Intel: ${DARWIN_AMD64_URL}"
echo "  ARM64: ${DARWIN_ARM64_URL}"
echo ""
echo "SHA256:"
echo "  Intel: ${SHA256_AMD64}"
echo "  ARM64: ${SHA256_ARM64}"
echo ""
echo "To publish to homebrew-sultengutt, copy this file to your tap repository"