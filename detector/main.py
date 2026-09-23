from concurrent import futures

import grpc

from pathlib import Path

import protobuf.compiled.inference_pb2_grpc
import protobuf.compiled.inference_pb2
import foundation.configuration

import foundation.model


model = foundation.model.new_model(Path("anti-spoof-mn3_float32.tflite"))
configuration = foundation.configuration.new_configuration(Path("config.json"))


class InferenceService(
    protobuf.compiled.inference_pb2_grpc.InferenceServiceServicer
):
    def Predict(self, request, context):
        if not request.image:
            context.set_code(
                grpc.StatusCode.INVALID_ARGUMENT
            )
            context.set_details("Image is empty")
            return protobuf.compiled.inference_pb2.PredictionResponse()

        try:
            result = model.execute(request.image)

            return protobuf.compiled.inference_pb2.PredictionResponse(
                result=result["result"],
                live=result["live"],
                spoof=result["spoof"],
            )

        except Exception as error:
            context.set_code(
                grpc.StatusCode.INTERNAL
            )
            context.set_details(str(error))
            return protobuf.compiled.inference_pb2.PredictionResponse()


def main():
    server = grpc.server(
        futures.ThreadPoolExecutor(max_workers=1)
    )

    protobuf.compiled.inference_pb2_grpc.add_InferenceServiceServicer_to_server(
        InferenceService(),
        server,
    )

    server.add_insecure_port(
        "localhost" + configuration.grpc.port
    )

    server.start()

    print("gRPC server running on " + configuration.grpc.port)

    server.wait_for_termination()


if __name__ == "__main__":
    main()
