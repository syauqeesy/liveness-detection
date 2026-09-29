import grpc

from service.prediction import (
    InvalidImageError,
    PredictionError,
    PredictionService,
)
from pb.compiled.inference_pb2 import (
    PredictionRequest,
    PredictionResponse,
)
from pb.compiled.inference_pb2_grpc import (
    InferenceServiceServicer,
)


class InferenceService(InferenceServiceServicer):
    def __init__(
        self,
        prediction_service: PredictionService,
    ):
        self._prediction_service = prediction_service

    def Predict(
        self,
        request: PredictionRequest,
        context: grpc.ServicerContext,
    ) -> PredictionResponse:
        try:
            result = self._prediction_service.predict(
                request_id=request.request_id,
                image_base64=request.image,
            )

        except InvalidImageError as error:
            context.abort(
                grpc.StatusCode.INVALID_ARGUMENT,
                str(error),
            )

        except PredictionError as error:
            context.abort(
                grpc.StatusCode.INTERNAL,
                str(error),
            )

        return PredictionResponse(
            result=result.result,
            live=result.live,
            spoof=result.spoof,
        )
