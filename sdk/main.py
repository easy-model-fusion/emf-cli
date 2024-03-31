# You can code freely in this file. It is the entry point of the application.
# Uncomment this example if you have installed stabilityai/sdxl-turbo
#
# Note: This file is the entry point of the application, it will be used as the main file to build the executable.
#
# from sdk.models import ModelsManagement
# from sdk import StabilityaiSdxlTurbo
# from sdk.options import Devices, OptionsTextToImage
#
# if __name__ == '__main__':
#     model_management = ModelsManagement()
#     model_stabilityai = StabilityaiSdxlTurbo()
#     options = OptionsTextToImage(
#             prompt="Astronaut in a jungle, cold color palette, "
#                    "muted colors, detailed, 8k",
#             device=Devices.GPU,
#             image_width=512,
#             image_height=512,
#     )
#     model_management.add_model(new_model=model_stabilityai, model_options=options)
#     model_management.load_model(StabilityaiSdxlTurbo.model_name)
#
#     image = model_management.generate_prompt()
#     image.show()
#        ).images[0]
#        image.show()

if __name__ == '__main__':
    print("Hello, EMF-World !")