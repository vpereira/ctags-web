FROM opensuse/leap:15.2

# for mongo-tools
RUN zypper ar https://download.opensuse.org/repositories/server:/database/openSUSE_Leap_15.2/server:database.repo
# RUN zypper ar https://download.opensuse.org/repositories/server:/database/openSUSE_Tumbleweed/server:database.repo
# universtal-ctags was dropped from leap 15.3 and tumbleweed
# RUN zypper ar https://download.opensuse.org/repositories/home:/vpereirabr/openSUSE_Tumbleweed/home:vpereirabr.repo
RUN zypper -n --gpg-auto-import-keys ref
RUN zypper -n in go screen git mongo-tools autoconf gcc make pkg-config libjansson-devel libxml2-devel libyaml-devel universal-ctags
