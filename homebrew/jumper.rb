require 'rbconfig'
class Jumper < Formula
  desc "This tool has been created for helping developers out with daily docker routines"
  homepage "https://github.com/asannikov/jumper"
  version "1.5.1"
  license "MIT"
  head "//github.com/asannikov/jumper.git"

  bottle :unneeded
  
  if Hardware::CPU.is_64_bit?
    case RbConfig::CONFIG["host_os"]
    when /mswin|msys|mingw|cygwin|bccwin|wince|emc/
      :windows
    when /darwin|mac os/
      url "https://github.com/asannikov/jumper/releases/download/v1.5.1/jumper-darwin-amd64.zip"
      sha256 "9ed11c20a484b2bf5e549d5ec8cd28c37cbba914ab1a17be8177a38e4b2b6eb1"
    when /linux/
      url "https://github.com/asannikov/jumper/releases/download/v1.5.1/jumper-linux-amd64.tar.gz"
      sha256 "7027c55c0481e6bf7e98c3cabfb7bcef3111fb269cb11eac6439181e5a5296b9"
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
      url "https://github.com/asannikov/jumper/releases/download/v1.5.1/jumper-linux-386.tar.gz"
      sha256 "3685e1822598c22bb2830ffd77d5c9e637f0e7a2a3268905016f1081b6553f98"
    when /solaris|bsd/
      :unix
    else
      :unknown
    end
  end

  system "brew tap tcnksm/ghr"
  depends_on "make" => :run
  depends_on "ghr" => :run
  depends_on 'go' => :build

  GOPATH = ENV["GOPATH"]
  def install
    ENV["GOPATH"] = GOPATH
    system "go get -u github.com/laher/goxc"
    bin.install "jumper"
  end

  test do
    system "jumper"
  end

end

