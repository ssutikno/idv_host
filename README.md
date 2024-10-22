# idv_host
Intelligent Desktop Virtualization (IDV) Host Management in UBUNTU

Install Required Components

sudo apt update
sudo apt install qemu-kvm libvirt-daemon-system libvirt-clients bridge-utils

go get github.com/libvirt/libvirt-go
go get github.com/google/uuid

# GIN web framework
go get github.com/gin-gonic/gin