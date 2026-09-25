from concurrent import futures
from signal import SIGINT, SIGTERM, signal

import grpc

from foundation.configuration import Configuration
from foundation.logger import Logger


class Grpc:
    def __init__(
        self,
        configuration: Configuration,
        logger: Logger,
    ):
        self._configuration = configuration
        self._logger = logger

        self._server = grpc.server(futures.ThreadPoolExecutor(max_workers=1))

        self._server.add_insecure_port("localhost" + configuration.grpc.port)

        signal(SIGTERM, self._shutdown)
        signal(SIGINT, self._shutdown)

    @property
    def server(self) -> grpc.Server:
        return self._server

    def start(self):
        self._server.start()

        self._logger.info("server running on " + self._configuration.grpc.port)

    def wait_for_termination(self):
        self._server.wait_for_termination()

    def _shutdown(self, signum, frame):
        self._logger.info("shutdown signal received")

        self._server.stop(grace=15)

        self._logger.info("application shutdown completed")


def new_grpc(configuration: Configuration, logger: Logger) -> Grpc:
    return Grpc(configuration, logger)
