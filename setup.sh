docker build -t madhush/ebpf-mac-golang .
docker image push madhush/ebpf-mac-golang:latest
kubectl apply -f depl.yaml
docker run -it --rm --privileged -p8000:8000 -v /lib/modules:/lib/modules:ro -v /etc/localtime:/etc/localtime:ro --pid=host madhush/ebpf-mac-golang
