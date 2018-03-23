FROM opensuse:42.3

RUN zypper ar https://download.opensuse.org/repositories/home:/vpereirabr/openSUSE_Leap_42.3/home:vpereirabr.repo
RUN zypper -n --gpg-auto-import-keys ref
RUN zypper -n up
RUN zypper -n in go screen git mongodb-shell
RUN zypper -n in autoconf gcc make pkg-config
RUN zypper -n in libjansson-devel libxml2-devel libyaml-devel
RUN zypper -n in universal-ctags
