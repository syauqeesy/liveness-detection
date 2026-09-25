import grpc

from foundation.configuration import Configuration
from foundation.logger import Logger
from foundation.model import AntiSpoofModel
from pb.compiled.inference_pb2_grpc import \
    add_InferenceServiceServicer_to_server
from service.inference import InferenceService


def new_service(
    server: grpc.Server,
    configuration: Configuration,
    logger: Logger,
    model: AntiSpoofModel,
):
    inference_service = InferenceService(configuration, logger, model)

    add_InferenceServiceServicer_to_server(inference_service, server)
