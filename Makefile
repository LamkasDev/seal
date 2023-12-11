SEALPLATFORM=seal64
SEALDIR=$(USERPROFILE)\Desktop\seal

GOOS=windows
GOARCH=amd64
GOTAGS=$(SEALPLATFORM),sealdebug

.PHONY: buildengine build install runengine run clean
buildengine:
	@set GOOS=$(GOOS)
	@set GOARCH=$(GOARCH)
	@go build -o build/$(SEALPLATFORM)/seal_engine.exe -tags $(GOTAGS) cmd/engine/main.go

build: buildengine

install: build
	@if exist "$(SEALDIR)\bin\$(SEALPLATFORM)" rmdir /S /Q "$(SEALDIR)\bin\$(SEALPLATFORM)"
	@xcopy "build\$(SEALPLATFORM)" "$(SEALDIR)\bin\$(SEALPLATFORM)\" /E /C /I >nul
	@if exist"$(SEALDIR)\resources" rmdir /S /Q "$(SEALDIR)\resources"
	@xcopy "resources" "$(SEALDIR)\resources\" /E /C /I >nul

runengine: buildengine
	@if not exist "$(SEALDIR)\bin\dev" mkdir "$(SEALDIR)\bin\dev"
	@copy "build\$(SEALPLATFORM)\seal.exe" "$(SEALDIR)\bin\dev\seal_engine.exe" >nul
	@cd "build\$(SEALPLATFORM)" && .\seal_engine.exe

run: build
	@if exist "$(SEALDIR)\bin\dev" rmdir /S /Q "$(SEALDIR)\bin\dev"
	@xcopy "build\$(SEALPLATFORM)" "$(SEALDIR)\bin\dev\" /E /C /I >nul
	@cd "build\$(SEALPLATFORM)" && .\seal_engine.exe

clean:
	@if exist "build" rmdir /S /Q build