require 'rbconfig'
class Jumper < Formula
  desc "This tool has been created for helping developers out with daily docker routines"
  homepage "https://github.com/asannikov/jumper"
  version "1.5.2"
  license "MIT"
  head "//github.com/asannikov/jumper.git"

  if Hardware::CPU.is_64_bit?
    case RbConfig::CONFIG["host_os"]
    when /mswin|msys|mingw|cygwin|bccwin|wince|emc/
      :windows
    when /darwin|mac os/
      url "https://github.com/asannikov/jumper/releases/download/v1.5.2/jumper-darwin-amd64.zip"
      sha256 "295bd0a026315c069b39cb2e859069e25e13ec121b17bd3441a852a839252b28"
    when /linux/
      url "https://github.com/asannikov/jumper/releases/download/v1.5.2/jumper-linux-amd64.tar.gz"
      sha256 "1a6cb860db0d38b6a609b7efa6c3c0e35223effd24702c2caad06d6c4c668143"
    when /solaris|bsd/
      :unix
    else
      :unknown
    end
  else
    case RbConfig::CONFIG["host_os"]
    when /mswin|msys|mingw|cygwin|bccwin|wince|emc/
      :windows
    when /linux/
      url "https://github.com/asannikov/jumper/releases/download/v1.5.2/jumper-linux-386.tar.gz"
      sha256 "1edb8dbb5c36b42e3dd4ee3502d36c3fbf4a8a02e0af6536a167ffd65544859e"
    when /solaris|bsd/
      :unix
    else
      :unknown
    end
  end

  def install
    bin.install "jumper"
  end

  test do
    system "jumper"
  end

end

