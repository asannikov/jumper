require "rbconfig"
class Jumper < Formula
  desc "Tool helps developers out with daily docker monotonous routine"
  homepage "https://github.com/asannikov/jumper"
  version "1.7.1"
  license "MIT"
  head "github.com/asannikov///jumper.git"

  if Hardware::CPU.is_64_bit?
    case RbConfig::CONFIG["host_os"]
    when /mswin|msys|mingw|cygwin|bccwin|wince|emc/
      :windows
    when /darwin|mac os/
      url "https://github.com/asannikov/jumper/releases/download/v1.7.1/jumper-darwin-amd64.zip"
      sha256 "4949c8465e7c9b7f107c225ef3c0a0bf92a02ca3b87eeabbe97ce9f011a3039f"
    when /linux/
      url "https://github.com/asannikov/jumper/releases/download/v1.7.1/jumper-linux-amd64.tar.gz"
      sha256 "f47c521f9842729336a8e7fb5610e8fb3dd8f2f0b2af1d2873f09b788f2161c0"
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
      url "https://github.com/asannikov/jumper/releases/download/v1.7.1/jumper-linux-386.tar.gz"
      sha256 "c01a6ae6299b5199080938d9a1bd1804c9dc0b5197ac786e5ffd4f173711883f"
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