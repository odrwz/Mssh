class Mssh < Formula
  desc "SSH Host & Key Manager CLI"
  homepage "https://github.com/odrwz/CLImssh"
  url "https://raw.githubusercontent.com/odrwz/CLImssh/main/climssh"
  # sha256 "<run: sha256sum climssh> and paste here after tagging"
  version "2.0.0"

  def install
    bin.install "climssh"
  end

  test do
    assert_match "CLImssh", shell_output("#{bin}/climssh 2>&1 || true")
  end
end
