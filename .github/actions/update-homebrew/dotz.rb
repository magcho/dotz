class Dotz < Formula
  desc "backup and restore dotfiles"
  homepage "https://github.com/magcho/dotz/"
  version "1.2.3"
  url "https://github.com/magcho/dotz/releases/download/v#{version}/dotz"
  sha256 "abcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcd"

  def install
    bin.install "dotz"
  end

  test do
    system "false"
  end
end
