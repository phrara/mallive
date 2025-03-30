.PHONY: gen
gen: genproto genopenapi


.PHONY: genproto
genproto:
	@./scripts/genproto.sh

.PHONY: genopenapi
genopenapi:
	@./scripts/genopenapi.sh

.PHONY: root
root:
	cd ~/go-proj/mallive