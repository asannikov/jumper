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
      sha256 "aea48d672bfe8b51352d14256c80cdba49855ec018e72094656dc1f228fdd1df"
    when /linux/
      url "https://github.com/asannikov/jumper/releases/download/v1.7.1/jumper-linux-amd64.tar.gz"
      sha256 "11a541cfd1c4284b9d5c60ad81d50155120db8f566e6aab23b127018ec80d85f"
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
      sha256 "f6e1ef07d2af0fdc68b72ca380782d3a190b8d013c03fdb8205d6c249a669702"
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