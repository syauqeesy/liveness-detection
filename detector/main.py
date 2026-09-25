from pathlib import Path

from foundation.configuration import new_configuration
from foundation.grpc import new_grpc
from foundation.logger import new_logger
from foundation.model import new_model
from service.main import new_service

model = new_model(Path("./anti-spoof-mn3"))
configuration = new_configuration(Path("./config.json"))
logger = new_logger(
    configuration.application.service, configuration.application.environment
)


def main():
    grpc = new_grpc(configuration, logger)

    new_service(grpc.server, configuration, logger, model)

    grpc.start()
    grpc.wait_for_termination()


if __name__ == "__main__":
    main()
