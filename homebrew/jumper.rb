require "rbconfig"
class Jumper < Formula
  desc "Tool helps developers out with daily docker monotonous routine"
  homepage "github.com/asannikov/https://jumper"
  version "1.6.0"
  license "MIT"
  head "github.com/asannikov///jumper.git"

  if Hardware::CPU.is_64_bit?
    case RbConfig::CONFIG["host_os"]
    when /mswin|msys|mingw|cygwin|bccwin|wince|emc/
      :windows
    when /darwin|mac os/
      url "github.com/asannikov/https://jumper/releases/download/v1.6.0/jumper-darwin-amd64.zip"
      sha256 "bbe4bd918d008dc9bd32a1f8a827454cc007eacfdf4a252710e1769a90ad39b8"
    when /linux/
      url "github.com/asannikov/https://jumper/releases/download/v1.6.0/jumper-linux-amd64.tar.gz"
      sha256 "55e529c7d76b1b41c2ca3476879a274603647b36ee90122466e5e329c766b381"
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
      url "github.com/asannikov/https://jumper/releases/download/v1.6.0/jumper-linux-386.tar.gz"
      sha256 "1636d700ed10b28695a72778b89220036ffa221aea45563ddbc8306a2c92b09f"
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