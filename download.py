import argparse
import logging
import sys

import diffusers
import transformers

# Configuration file keys
MODELS_KEY = 'models'
MODEL_NAME_KEY = 'name'
MODEL_CONFIG_KEY = 'config'
CONFIG_MODULE_KEY = 'modulename'
CONFIG_CLASS_KEY = 'classname'

# Authorized library names for download
AUTHORIZED_MODULE_NAMES = {'diffusers', 'transformers'}

# Map module name to a default class name
model_config_default_class_for_module = {
    'diffusers': "DiffusionPipeline",
    'transformers': "AutoModel",
}

# Set up logging
logging.basicConfig(level=logging.INFO)


def download(downloads_path, model_name, module_name, class_name):
    """
    Downloads the model given data.

    Args:
        downloads_path (str): Path to the download's directory.
        model_name (str): The name of the model to download.
        module_name (str): The name of the module to use for downloading.
        class_name (str): The class within the module to use for downloading.

    Returns:
        None
    """

    def print_error(message):
        print(message, file=sys.stderr)

    # Check if the model name is not empty
    if model_name is None or model_name.strip() == '':
        print_error(f"ERROR: Model '{model_name}' is invalid.")
        sys.exit(1)

    # Check if the module is authorized
    if module_name not in AUTHORIZED_MODULE_NAMES:
        print_error(f"Module '{module_name}' is not authorized.")
        sys.exit(1)

    # If class is not provided, use the default
    if class_name is None or class_name.strip() == '':
        logging.warning("Module class not provided, using default but might fail.")
        class_name = model_config_default_class_for_module.get(module_name)

    try:
        # Transforming from strings to actual objects
        module_obj = globals()[module_name]
        class_obj = getattr(module_obj, class_name)

        # TODO : check if the downloads_path exists

        # Downloading the model
        model = class_obj.from_pretrained(model_name)
        model.save_pretrained(downloads_path + model_name)
        # TODO : Tokenizer?
        # TODO : Options?

    except Exception as e:
        print_error(f"Error while downloading model {model_name}: {e}")
        sys.exit(1)


def parse_arguments():
    """
    Parse command-line arguments.

    Returns:
        argparse.Namespace: Parsed command-line arguments.
    """

    parser = argparse.ArgumentParser(description="Script to download a specific model.")

    parser.add_argument("downloads_path", type=str, help="Path to the downloads directory")
    parser.add_argument("model_name", type=str, help="Model name")
    parser.add_argument("module_name", type=str, help="Module name")
    parser.add_argument("class_name", nargs="?", type=str, help="Class name (optional)")

    return parser.parse_args()


def main():
    """
    Main function to execute the download process based on the provided configuration file.
    """

    args = parse_arguments()

    # Run download with specified arguments
    download(args.downloads_path, args.model_name, args.module_name, args.class_name)


if __name__ == "__main__":
    main()
