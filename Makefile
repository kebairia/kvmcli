# --- Macros ---
FILES = {config,src,kvmcli}
TARGET = ~/.local/bin/kvmcli
# --- Targets ---
clean:
	@echo "Cleaning..."
	rm -rf ${TARGET}

install: 
	mkdir -p ${TARGET}
	cp -arv ${FILES} ${TARGET}
	find ${TARGET} -iname "__pycache__" -type d -exec rm -rf {} \;


.PHONY: clean
