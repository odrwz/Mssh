class Mssh < Formula
  desc "SSH Host & Key Manager CLI"
  homepage "https://github.com/odrwz/mssh"
  url "https://raw.githubusercontent.com/odrwz/mssh/main/mssh"
  sha256 "beae60ddfc4ec21c21f2929b0eaa3547c9dc2cb9057c0e8f4be6820ef2a0ad71"
  version "2.0.0"

  def install
    bin.install "mssh"
  end

  test do
    assert_match "mssh", shell_output("#{bin}/mssh 2>&1 || true")
  end
end
