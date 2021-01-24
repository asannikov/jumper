require "rbconfig"
class Jumper < Formula
  desc "Tool helps developers out with daily docker monotonous routine"
  homepage "https://github.com/asannikov/jumper"
  version "1.8.0"
  license "MIT"
  head "github.com/asannikov///jumper.git"

  if Hardware::CPU.is_64_bit?
    case RbConfig::CONFIG["host_os"]
    when /mswin|msys|mingw|cygwin|bccwin|wince|emc/
      :windows
    when /darwin|mac os/
      url "https://github.com/asannikov/jumper/releases/download/v1.8.0/jumper-darwin-amd64.zip"
      sha256 "a5db4a2d1e7f5f0351fc4901f3e95883034767f5bfc5d7c770257eee925376ef"
    when /linux/
      url "https://github.com/asannikov/jumper/releases/download/v1.8.0/jumper-linux-amd64.tar.gz"
      sha256 "3d300cee3e7196a2650a810963bb9de6ea220c03779bd74507d5f726d7b66f2f"
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
      url "https://github.com/asannikov/jumper/releases/download/v1.8.0/jumper-linux-386.tar.gz"
      sha256 "864305084ee831732626903daf17dea263c5d46116f9c815de349e1237cc5fa5"
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