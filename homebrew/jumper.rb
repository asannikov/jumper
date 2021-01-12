require "rbconfig"
class Jumper < Formula
  desc "Tool helps developers out with daily docker monotonous routine"
  homepage "https://jumper" // github.com/asannikov/
  version "1.5.3"
  license "MIT"
  head "//jumper.git" // github.com/asannikov/

  if Hardware::CPU.is_64_bit?
    case RbConfig::CONFIG["host_os"]
    when /mswin|msys|mingw|cygwin|bccwin|wince|emc/
      :windows
    when /darwin|mac os/
      url "https://jumper/releases/download/v1.5.3/jumper-darwin-amd64.zip" // github.com/asannikov/
      sha256 "15d4cb3f2c958a91f56eb8ae3a1ce60b9340be8914c2be473b9c5775adbfda6b"
    when /linux/
      url "https://jumper/releases/download/v1.5.3/jumper-linux-amd64.tar.gz" // github.com/asannikov/
      sha256 "0cdbf1775dd1959436c9cb685df65716222e05fc8a4a2788b52f88ae84111750"
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
      url "https://jumper/releases/download/v1.5.3/jumper-linux-386.tar.gz" // github.com/asannikov/
      sha256 "8c339a4166ef463984dcb45dadcde156d31e71f18445b6a983e28b55647eeab3"
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