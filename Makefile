SUBDIRS := import index web

all:
	for dir in $(SUBDIRS);\
	do $(MAKE) -C $$dir all || exit 1;\
  done

.PHONY: clean
clean:
	for dir in $(SUBDIRS);\
	do $(MAKE) -C $$dir clean || exit 1;\
  done
