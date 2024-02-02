import yaml
import logging
import diffusers
import transformers

# Paths
CONFIGURATION_FILE_PATH = 'config.yaml'
DOWNLOADED_MODELS_PATH = 'models/'

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


def download(model_name, module_name, class_name):
    """
    Downloads the model given data.

    Args:
        model_name (str): The name of the model to download.
        module_name (str): The name of the module to use for downloading.
        class_name (str): The class within the module to use for downloading.

    Returns:
        None
    """

    # TODO : check that the model does exist using the API?
    # Check if the model name is not empty
    if model_name is None or model_name.strip() == '':
        logging.error(f"Model '{model_name}' is invalid.")
        return

    # TODO : sdk or client : if model has already been downloaded = ask to overwrite or cancel?

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

        # Downloading the model
        model = class_obj.from_pretrained(model_name)
        model.save_pretrained(DOWNLOADED_MODELS_PATH + model_name)
        # TODO : Tokenizer?
        # TODO : Options?

        logging.info(f"Model {model_name} saved.")
        # TODO : update the model path in the config? client or sdk?

    except Exception as e:
        logging.error(f"Error while downloading model {model_name}: {e}")


def load_config():
    """
    Loads YAML configuration data from a file.

    Returns:
        dict: The loaded YAML data.
    """

    try:
        with open(CONFIGURATION_FILE_PATH, 'r') as file:
            data = yaml.safe_load(file)
        return data
    except FileNotFoundError:
        logging.error(f"File not found: {CONFIGURATION_FILE_PATH}")
    except yaml.YAMLError as e:
        logging.error(f"Error reading YAML file {CONFIGURATION_FILE_PATH}: {e}")
    except Exception as e:
        logging.error(f"An unexpected error occurred: {e}")


def download_models(models):
    """
    Downloads multiple models based on the provided item.

    Args:
        models (list): List of dictionaries containing models.

    Returns:
        None
    """

    # Models key not found in the provided item
    if not models:
        logging.warning(f"No '{MODELS_KEY}' key found in the provided item.")
        return

    # Models value is not a dictionary
    if not isinstance(models, list):
        logging.warning(f"'{MODELS_KEY}' should be a list in the provided item.")
        return

    # Download every model
    for model in models:
        download_model(model)


def download_model(model):
    """
    Downloads a single model based on the provided item.

    Args:
        model (dict): Dictionary containing model.

    Returns:
        None
    """

    # Name key not found in the model
    if MODEL_NAME_KEY not in model:
        logging.warning(f"'{MODEL_NAME_KEY}' key not found in the model.")
        return None, None

    # Extracting the model configuration datas necessary to download
    module_name, class_name = get_model_config_datas(model)

    # Actually downloading the model
    download(model[MODEL_NAME_KEY], module_name, class_name)


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
        logging.warning(f"Each item under '{MODELS_KEY}' should be a dictionary.")
        return None, None

    # Configuration key not found in the model
    if MODEL_CONFIG_KEY not in model:
        logging.warning(f"'{MODEL_CONFIG_KEY}' key not found in the model.")
        return None, None

    # Configuration key not found in the model
    config = model[MODEL_CONFIG_KEY]
    if not isinstance(config, dict):
        logging.warning(f"The '{MODEL_CONFIG_KEY}' key should be a dictionary in the model.")
        return None, None

    # Returning library name and class
    return config.get(CONFIG_MODULE_KEY), config.get(CONFIG_CLASS_KEY)


def main():
    """
    Main function to execute the download process based on the provided configuration file.
    """

    config_data = load_config()

    if not config_data:
        logging.error(f"Failed to load YAML data from {CONFIGURATION_FILE_PATH}")
        return

    download_models(config_data.get(MODELS_KEY, []))


if __name__ == "__main__":
    main()
