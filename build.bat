go build main.go
if exist "shavis.exe" del "shavis.exe"
rename "main.exe" "shavis.exe"