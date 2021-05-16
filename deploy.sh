# Build go
env GOOS=linux GOARCH=amd64 go build -o droppy_cgnorthpeak_go *.go

# Stop service
ssh droppy "service go-droppy-cgnorthpeak stop"

# Remove File
ssh droppy "rm /var/www/droppy_cgnorthpeak_go/droppy_cgnorthpeak_go"

# Upload file
scp droppy_cgnorthpeak_go droppy:/var/www/droppy_cgnorthpeak_go

# Chmod
ssh droppy "chmod +x /var/www/droppy_cgnorthpeak_go/droppy_cgnorthpeak_go"

# Stop service
ssh droppy "service go-droppy-cgnorthpeak stop"

# Refresh daemons
ssh droppy "systemctl daemon-reload"

# Start service again
ssh droppy "service go-droppy-cgnorthpeak start"
