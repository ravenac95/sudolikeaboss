package:
	-rm -rf _build
	CC=clang go build
	mkdir -p _build/amd64
	mv sudolikeaboss _build/amd64/sudolikeaboss
	@cd _build/amd64; zip -r sudolikeaboss_`cat ../../VERSION`_darwin_amd64.zip .
