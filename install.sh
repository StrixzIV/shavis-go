go install 

if [ ! -e ~/.shavis-go.yaml ]; then
echo "Copying config file..."
cp .shavis-go.yaml ~/
fi

echo "Installation complete!"
echo "shavis-go config are stored in $HOME/.shavis-go.yaml"