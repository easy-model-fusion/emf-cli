import argparse
import logging
import os
import sys

import diffusers
import transformers

DIFFUSERS = 'diffusers'
TRANSFORMERS = 'transformers'

# Authorized library names for download
AUTHORIZED_MODULE_NAMES = {DIFFUSERS, TRANSFORMERS}

# Map module name to a default class name
model_config_default_class_for_module = {
    DIFFUSERS: "DiffusionPipeline",
    TRANSFORMERS: "AutoModel",
}

# Set up logging
logging.basicConfig(level=logging.INFO)


def download(downloads_path, model_name, module_name, class_name, overwrite=False):
    """
    Downloads the model given data.

    Args:
        downloads_path (str): Path to the download's directory.
        model_name (str): The name of the model to download.
        module_name (str): The name of the module to use for downloading.
        class_name (str): The class within the module to use for downloading.
        overwrite (bool): Whether to overwrite the downloaded model if it exists.

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

    # Class is not provided, trying the default one
    if class_name is None or class_name.strip() == '':
        logging.warning("Module class not provided, using default but might fail.")
        class_name = model_config_default_class_for_module.get(module_name)

    try:
        # Transforming from strings to actual objects
        module_obj = globals()[module_name]
        class_obj = getattr(module_obj, class_name)

        # Local path where the model will be downloaded
        model_path = os.path.join(downloads_path, model_name)

        # Check if the model_path already exists
        if not overwrite and os.path.exists(model_path):
            print_error(f"Directory '{model_path}' already exists.")
            sys.exit(1)

        # Downloading the model
        # TODO : Options?
        model = class_obj.from_pretrained(model_name)
        model.save_pretrained(model_path)

        # TODO : Tokenizer?
        if module_name == TRANSFORMERS:
            tokenizer_path = os.path.join(model_path, 'tokenizer')
            tokenizer = transformers.AutoTokenizer.from_pretrained(model_name)
            tokenizer.save_pretrained(tokenizer_path)

    except Exception as e:
        print_error(f"Error while downloading model {model_name}: {e}")
        sys.exit(2)


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
    parser.add_argument("--overwrite", action="store_true", help="Overwrite existing directories", dest="overwrite")

    return parser.parse_args()


def main():
    """
    Main function to execute the download process based on the provided configuration file.
    """

    args = parse_arguments()

    print(args.downloads_path, args.model_name, args.module_name, args.class_name, args.overwrite)

    # Run download with specified arguments
    download(args.downloads_path, args.model_name, args.module_name, args.class_name, args.overwrite)


if __name__ == "__main__":
    main()
