import grpc

from foundation.configuration import Configuration
from foundation.logger import Logger
from foundation.model import AntiSpoofModel
from pb.compiled.inference_pb2_grpc import add_InferenceServiceServicer_to_server
from service.grpc.inference import InferenceService
from service.prediction import PredictionService


def new_service(
    server: grpc.Server,
    logger: Logger,
    model: AntiSpoofModel,
):
    prediction_service = PredictionService(
        logger,
        model,
    )

    inference_service = InferenceService(
        prediction_service,
    )

    add_InferenceServiceServicer_to_server(
        inference_service,
        server,
    )
