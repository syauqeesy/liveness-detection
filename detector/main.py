from concurrent import futures
from pathlib import Path
from signal import SIGINT, SIGTERM, signal

import grpc

from foundation.configuration import new_configuration
from foundation.logger import new_logger
from foundation.model import new_model
from service.main import new_service

model = new_model(Path("./anti-spoof-mn3"))
configuration = new_configuration(Path("config.json"))
logger = new_logger(
    configuration.application.service, configuration.application.environment
)


def main():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=1))

    new_service(server, configuration, logger, model)

    server.add_insecure_port("localhost" + configuration.grpc.port)

    def shutdown(signum, frame):
        logger.info("shutdown signal received")

        server.stop(grace=15)

        logger.info("application shutdown completed")

    signal(SIGTERM, shutdown)
    signal(SIGINT, shutdown)

    server.start()

    logger.info("server running on " + configuration.grpc.port)

    server.wait_for_termination()


if __name__ == "__main__":
    main()
