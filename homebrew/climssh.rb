class Climssh < Formula
  desc "Interactive CLI tool to manage SSH config and keys"
  homepage "https://github.com/odrwz/CLImssh"
  head "https://github.com/odrwz/CLImssh.git", branch: "main"

  depends_on "go" => :build

  def install
    system "go", "build", "-o", bin/"climssh", "."
  end

  test do
    assert_match "CLImssh", shell_output("#{bin}/climssh 2>&1 || true")
  end
end
