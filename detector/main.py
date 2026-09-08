from pathlib import Path

from keras.applications import MobileNetV2
from keras.utils import get_file

from foundation.configuration import new_configuration


keras_dir = Path("./.keras")
model_dir = Path("./.keras/models")


def main():
    configuration = new_configuration(Path("./config.json"))

    model_dir.mkdir(parents=True, exist_ok=True)

    weights_path = get_file(
        fname=configuration.models.name,
        origin=configuration.models.url,
        cache_dir=keras_dir,
        cache_subdir="models",
    )

    model = MobileNetV2(weights=weights_path)

    model.summary()


if __name__ == "__main__":
    main()
