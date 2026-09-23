import json
from dataclasses import dataclass
from pathlib import Path


@dataclass
class ApplicationConfiguration:
    service: str
    environment: str


@dataclass
class GRPCConfiguration:
    port: str


@dataclass
class Configuration:
    application: ApplicationConfiguration
    grpc: GRPCConfiguration


def new_configuration(path: Path) -> Configuration:
    with path.open("r", encoding="utf-8") as file:
        data = json.load(file)

    return Configuration(
        application=ApplicationConfiguration(
            service=data["application"]["service"],
            environment=data["application"]["environment"],
        ),
        grpc=GRPCConfiguration(
            port=data["grpc"]["port"],
        ),
    )
