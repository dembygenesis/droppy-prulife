# Build go
env GOOS=linux GOARCH=amd64 go build -o droppy_prulife_go *.go

# Stop service
ssh droppy "service go-droppy-prulife stop"

# Remove File
ssh droppy "rm /var/www/droppy_prulife_go/droppy_prulife_go"

# Upload file
scp droppy_prulife_go droppy:/var/www/droppy_prulife_go

# Chmod
ssh droppy "chmod +x /var/www/droppy_prulife_go/droppy_prulife_go"

# Stop service
ssh droppy "service go-droppy-prulife stop"

# Refresh daemons
ssh droppy "systemctl daemon-reload"

# Start service again
ssh droppy "service go-droppy-prulife start"
