import json
from dataclasses import dataclass
from pathlib import Path


@dataclass
class ModelConfiguration:
    name: str
    url: str


@dataclass
class Configuration:
    models: ModelConfiguration


def new_configuration(path: Path) -> Configuration:
    with path.open("r", encoding="utf-8") as file:
        data = json.load(file)

    return Configuration(
        models=ModelConfiguration(
            name=data["models"]["name"],
            url=data["models"]["url"],
        ),

    )
