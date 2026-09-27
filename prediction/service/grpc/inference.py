import base64
import binascii
import time

import grpc
from foundation.configuration import Configuration
from foundation.logger import Logger
from foundation.model import AntiSpoofModel
from pb.compiled.inference_pb2 import PredictionRequest, PredictionResponse
from pb.compiled.inference_pb2_grpc import InferenceServiceServicer


class InferenceService(InferenceServiceServicer):
    def __init__(
        self,
        configuration: Configuration,
        logger: Logger,
        model: AntiSpoofModel,
    ):
        self._configuration = configuration
        self._logger = logger
        self._model = model

    def Predict(
        self,
        request: PredictionRequest,
        context: grpc.ServicerContext,
    ) -> PredictionResponse:
        try:
            image = base64.b64decode(
                request.image,
                validate=True,
            )

        except (binascii.Error, ValueError):
            self._logger.warn(
                "prediction rejected",
                "reason",
                "invalid base64 image",
            )

            context.abort(
                grpc.StatusCode.INVALID_ARGUMENT,
                "Image is not valid base64",
            )

        if not image:
            self._logger.warn(
                "prediction rejected",
                "reason",
                "empty image",
            )

            context.abort(
                grpc.StatusCode.INVALID_ARGUMENT,
                "Image is empty",
            )

        started = time.perf_counter()

        try:
            result = self._model.execute(image)

            return PredictionResponse(
                result=result.result,
                live=result.live,
                spoof=result.spoof,
            )

        except Exception as error:
            self._logger.error(
                "prediction failed",
                "error",
                str(error),
            )

            context.abort(
                grpc.StatusCode.INTERNAL,
                "Prediction failed",
            )

        finally:
            inference_time_ms = (
                time.perf_counter() - started
            ) * 1000

            self._logger.info(
                "prediction completed",
                "duration_ms",
                round(inference_time_ms, 2),
            )
