class Mssh < Formula
  desc "SSH Host & Key Manager CLI"
  homepage "https://github.com/odrwz/CLImssh"
  url "https://raw.githubusercontent.com/odrwz/CLImssh/main/climssh"
  sha256 "1dfa6ddd29c441fc242a109f36a49dbcfffcaf92a5714f289cb7163f18b23849"
  version "2.0.0"

  def install
    bin.install "climssh"
  end

  test do
    assert_match "CLImssh", shell_output("#{bin}/climssh 2>&1 || true")
  end
end
