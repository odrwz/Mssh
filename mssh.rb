class Mssh < Formula
  desc "SSH Host & Key Manager CLI"
  homepage "https://github.com/odrwz/CLImssh"
  url "https://raw.githubusercontent.com/odrwz/CLImssh/main/climssh"
  sha256 "4c911f120e53b4bfaae0ce9dde0729ff95e18fa283aa0b76b56c3ac86809135b"
  version "2.0.0"

  def install
    bin.install "climssh"
  end

  test do
    assert_match "CLImssh", shell_output("#{bin}/climssh 2>&1 || true")
  end
end
