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
