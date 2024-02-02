import argparse
import yaml
import logging
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

    # Check if the model name is not empty
    if model_name is None or model_name.strip() == '':
        logging.error(f"Model '{model_name}' is invalid.")
        return

    # Check if the module is authorized
    if module_name not in AUTHORIZED_MODULE_NAMES:
        logging.error(f"Module '{module_name}' is not authorized.")
        return

    # If class is not provided, use the default
    if class_name is None:
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

        logging.info(f"Model {model_name} saved.")

    except Exception as e:
        logging.error(f"Error while downloading model {model_name}: {e}")


def load_config(config_file_path):
    """
    Loads YAML configuration data from a file.

    Returns:
        dict: The loaded YAML data.
    """

    try:
        with open(config_file_path, 'r') as file:
            data = yaml.safe_load(file)
        return data
    except FileNotFoundError:
        logging.error(f"File not found: {config_file_path}")
    except yaml.YAMLError as e:
        logging.error(f"Error reading YAML file {config_file_path}: {e}")
    except Exception as e:
        logging.error(f"An unexpected error occurred: {e}")


def download_models(downloads_path, models):
    """
    Downloads multiple models based on the provided item.

    Args:
        downloads_path (str): Path to the download's directory.
        models (list): List of dictionaries containing models.

    Returns:
        None
    """

    # Models key not found in the provided item
    if not models:
        logging.error(f"No '{MODELS_KEY}' key found in the provided item.")
        return

    # Models value is not a dictionary
    if not isinstance(models, list):
        logging.error(f"'{MODELS_KEY}' should be a list in the provided item.")
        return

    # Download every model
    for model in models:
        download_model(downloads_path, model)


def download_model(downloads_path, model):
    """
    Downloads a single model based on the provided item.

    Args:
        downloads_path (str): Path to the download's directory.
        model (dict): Dictionary containing model.

    Returns:
        None
    """

    # Name key not found in the model
    if MODEL_NAME_KEY not in model:
        logging.error(f"'{MODEL_NAME_KEY}' key not found in the model.")
        return None, None

    # Extracting the model configuration datas necessary to download
    module_name, class_name = get_model_config_datas(model)

    # TODO : check that the model_name extracted from the config file does exist using the API?
    model_name = model[MODEL_NAME_KEY]

    # Actually downloading the model
    download(downloads_path, model_name, module_name, class_name)


def get_model_config_datas(model):
    """
    Extracts model specific configuration datas.

    Args:
        model (dict): Dictionary containing model.

    Returns:
        Tuple[str, str]: Module and class names.
    """

    # Model value is not a dictionary
    if not isinstance(model, dict):
        logging.error(f"Each item under '{MODELS_KEY}' should be a dictionary.")
        return None, None

    # Configuration key not found in the model
    if MODEL_CONFIG_KEY not in model:
        logging.error(f"'{MODEL_CONFIG_KEY}' key not found in the model.")
        return None, None

    # Configuration key not found in the model
    config = model[MODEL_CONFIG_KEY]
    if not isinstance(config, dict):
        logging.error(f"The '{MODEL_CONFIG_KEY}' key should be a dictionary in the model.")
        return None, None

    # Returning library name and class
    return config.get(CONFIG_MODULE_KEY), config.get(CONFIG_CLASS_KEY)


def parse_arguments():
    """
    Parse command-line arguments.

    Returns:
        argparse.Namespace: Parsed command-line arguments.
    """

    parser = argparse.ArgumentParser(description="Script to download models based on configuration or specific model.")

    subparsers = parser.add_subparsers(dest="command")

    # Subcommand for downloading all models from a configuration file
    config_parser = subparsers.add_parser("all", help="Download all models from a config file")
    config_parser.add_argument("downloads_path", type=str, help="Path to the downloads directory")
    config_parser.add_argument("file_path", type=str, help="Path to the config file")

    # Subcommand for downloading a specific model with its name, module and optional class
    model_parser = subparsers.add_parser("model", help="Download a specific model")
    model_parser.add_argument("downloads_path", type=str, help="Path to the downloads directory")
    model_parser.add_argument("model_name", type=str, help="Model name")
    model_parser.add_argument("module_name", type=str, help="Module name")
    model_parser.add_argument("class_name", nargs="?", type=str, help="Class name (optional)")

    return parser.parse_args()


def main():
    """
    Main function to execute the download process based on the provided configuration file.
    """

    args = parse_arguments()

    if args.command == 'all':
        # Load and process configuration file
        config_data = load_config(args.file_path)

        if not config_data:
            logging.error(f"Failed to load YAML data from {args.file_path}")
            return

        download_models(args.downloads_path, config_data.get(MODELS_KEY, []))
    elif args.command == 'model':
        # Run download with specified arguments
        download(args.downloads_path, args.model_name, args.module_name, args.class_name)


if __name__ == "__main__":
    main()
