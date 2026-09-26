# Liveness Detection
Setup protobuf tools:
```sh
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

sudo apt install protobuf-compiler
```

Compile protobuf:
```sh
# Prediction service
python3 -m grpc_tools.protoc -I pb --python_out=prediction/pb/compiled --pyi_out=prediction/protobuf/compiled --grpc_python_out=prediction/pb/compiled pb/inference.proto

# Backend service
protoc --proto_path=. --go_out=pb/compiled --go-grpc_out=pb/compiled pb/*.proto
```

Use CUDA inside python venv:
```sh
export LD_LIBRARY_PATH="$(find "$VIRTUAL_ENV/lib/python3.13/site-packages/nvidia" -type d -name lib -print | paste -sd:):${LD_LIBRARY_PATH}"
```

Run main.py:
```sh
PYTHONPATH="$PWD/pb/compiled:$PYTHONPATH" python3 main.py
```

Run main.go:
```sh
go run main.go http
```

Container setup:
```sh
docker build -t liveness:v1.0.0 .
docker run -d --rm --name liveness -p 5173:80 liveness:v1.0.0

docker build -t api-liveness:v1.0.0 .
docker run -d --rm --name api-liveness --add-host=host.docker.internal:host-gateway -p 3301:80 -v"$(pwd)/config.json:/app/config.json:ro" api-liveness:v1.0.0

docker build -t prediction-liveness:v1.0.0 .
docker run -d --gpus all --rm --name prediction-liveness -p 3302:80 -v"$(pwd)/config.json:/app/config.json:ro" prediction-liveness:v1.0.0
```

Setup nvidia container toolkit:
```sh
sudo apt-get install ca-certificates curl gnupg2

curl -fsSL https://nvidia.github.io/libnvidia-container/gpgkey | sudo gpg --dearmor -o /usr/share/keyrings/nvidia-container-toolkit-keyring.gpg
curl -s -L https://nvidia.github.io/libnvidia-container/stable/deb/nvidia-container-toolkit.list | sed 's#deb https://#deb [signed-by=/usr/share/keyrings/nvidia-container-toolkit-keyring.gpg] https://#g' | sudo tee /etc/apt/sources.list.d/nvidia-container-toolkit.list

sudo apt update

sudo apt-get install nvidia-container-toolkit

sudo nvidia-ctk runtime configure --runtime=docker
```
