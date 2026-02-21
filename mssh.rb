class Mssh < Formula
  desc "SSH Host & Key Manager CLI"
  homepage "https://github.com/odrwz/mssh"
  url "https://raw.githubusercontent.com/odrwz/mssh/main/mssh"
  sha256 "38c68cc12af17b0c8dd558a337ebd1f7ca9d71a92e357b1e4d47a327816b4ba1"
  version "2.0.0"

  def install
    bin.install "mssh"
  end

  test do
    assert_match "mssh", shell_output("#{bin}/mssh 2>&1 || true")
  end
end
