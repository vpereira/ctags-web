FROM opensuse/tumbleweed:latest

RUN zypper -n in go screen git libjansson-devel libxml2-devel libyaml-devel universal-ctags
