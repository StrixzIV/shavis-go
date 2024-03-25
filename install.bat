@echo off
go install

echo Copying config file...
copy .shavis-go.yaml %userprofile%

echo Installation complete!
echo shavis-go config are stored in "%userprofile%\.shavis-go.yaml"