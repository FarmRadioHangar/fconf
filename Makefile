VERSION=0.1.0
NAME=fconf_$(VERSION)
OUT_DIR=bin/linux_arm/fconf_$(VERSION)

all:$(OUT_DIR)/fconf
$(OUT_DIR)/fconf:main.go
	gox  \
		-output "bin/{{.OS}}_{{.Arch}}/{{.Dir}}_$(VERSION)/{{.Dir}}" \
		-osarch "linux/arm" github.com/FarmRadioHangar/fconf

tar:
	tar -zcvf fconf_$(VERSION).tar.gz  $(OUT_DIR)