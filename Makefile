.PHONY: build install clean test gui cli install-gui

BINARY_NAME=bill
INSTALL_PATH=$(GOPATH)/bin
BUILD_DIR=build
APP_NAME=Bill.app
APP_DIR=$(BUILD_DIR)/$(APP_NAME)
APP_CONTENTS=$(APP_DIR)/Contents
APP_MACOS=$(APP_CONTENTS)/MacOS
APP_RESOURCES=$(APP_CONTENTS)/Resources

cli:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/bill

gui:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME)-gui ./cmd/desktop

build: cli gui

install: build
	cp $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_PATH)/$(BINARY_NAME)
	cp $(BUILD_DIR)/$(BINARY_NAME)-gui $(INSTALL_PATH)/$(BINARY_NAME)-gui

install-gui: gui
	mkdir -p $(APP_MACOS)
	mkdir -p $(APP_RESOURCES)
	cp $(BUILD_DIR)/$(BINARY_NAME)-gui $(APP_MACOS)/$(BINARY_NAME)
	echo '<?xml version="1.0" encoding="UTF-8"?>' > $(APP_CONTENTS)/Info.plist
	echo '<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">' >> $(APP_CONTENTS)/Info.plist
	echo '<plist version="1.0">' >> $(APP_CONTENTS)/Info.plist
	echo '<dict>' >> $(APP_CONTENTS)/Info.plist
	echo '    <key>CFBundleExecutable</key>' >> $(APP_CONTENTS)/Info.plist
	echo '    <string>$(BINARY_NAME)</string>' >> $(APP_CONTENTS)/Info.plist
	echo '    <key>CFBundleIdentifier</key>' >> $(APP_CONTENTS)/Info.plist
	echo '    <string>com.louis.invoice-generator</string>' >> $(APP_CONTENTS)/Info.plist
	echo '    <key>CFBundleName</key>' >> $(APP_CONTENTS)/Info.plist
	echo '    <string>Bill</string>' >> $(APP_CONTENTS)/Info.plist
	echo '    <key>CFBundlePackageType</key>' >> $(APP_CONTENTS)/Info.plist
	echo '    <string>APPL</string>' >> $(APP_CONTENTS)/Info.plist
	echo '    <key>CFBundleShortVersionString</key>' >> $(APP_CONTENTS)/Info.plist
	echo '    <string>1.0</string>' >> $(APP_CONTENTS)/Info.plist
	echo '    <key>LSMinimumSystemVersion</key>' >> $(APP_CONTENTS)/Info.plist
	echo '    <string>10.13</string>' >> $(APP_CONTENTS)/Info.plist
	echo '</dict>' >> $(APP_CONTENTS)/Info.plist
	echo '</plist>' >> $(APP_CONTENTS)/Info.plist
	@echo "Created $(APP_NAME) in $(BUILD_DIR)"
	@echo "To install, run: cp -r $(APP_DIR) /Applications"

clean:
	rm -rf $(BUILD_DIR)
	go clean

test:
	go test -v ./...

run: gui
	./$(BUILD_DIR)/$(BINARY_NAME)-gui