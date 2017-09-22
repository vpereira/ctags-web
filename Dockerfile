FROM opensuse:42.3

RUN zypper -n ref
RUN zypper -n up
RUN zypper -n in go screen mongodb
RUN zypper -n in autoconf gcc make pkg-config
RUN zypper -n in libjansson-devel libxml2-devel libyaml-devel

RUN mkdir -p /data/db
# RUN mongod --nojournal --smallfiles --fork --syslog
