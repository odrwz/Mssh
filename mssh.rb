class Mssh < Formula
  desc "Interactive CLI tool to manage SSH Configs and Keys"
  homepage "https://github.com/qwer1234/climssh"
  # This URL will be replaced by the actual release tarball URL once published
  url "https://github.com/qwer1234/climssh/archive/refs/tags/v1.0.0.tar.gz"
  version "1.0.0"
  sha256 "0000000000000000000000000000000000000000000000000000000000000000" # Update this
  license "MIT"

  depends_on "go" => :build

  def install
    # Build the binary
    system "go", "build", *std_go_args(output: bin/"mssh"), "main.go", "sshconfig.go", "sshkey.go", "ui.go"
  end

  test do
    # Simple test to verify the binary runs
    assert_match "Welcome to CLImssh", shell_output("#{bin}/mssh --help", 1) # Because we don't handle --help properly, it might just exit or crash, but we can refine test later
  end
end
