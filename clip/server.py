
from concurrent import futures
import threading
from PIL import Image
from transformers import CLIPModel, CLIPProcessor, CLIPConfig
import clip_pb2
import clip_pb2_grpc
import grpc
import io
import os
import torch

MODEL_ID = "openai/clip-vit-base-patch32"

lock = threading.Lock()

class EncodingServer(clip_pb2_grpc.EncodingService):
    def __init__(self) -> None:
        self.device = "cuda" if torch.cuda.is_available() else "cpu"
        self.model = CLIPModel.from_pretrained(
            MODEL_ID,
        ).to(self.device)

        self.processor = CLIPProcessor.from_pretrained(MODEL_ID)
        self.config = CLIPConfig.from_pretrained(MODEL_ID)

        self.logit_scale = self.model.logit_scale.item(
        ) if self.model.logit_scale.item() else 4.60517
        
        print("Model clip loaded", "device:", self.device)

    def EncodeText(self, request, context):
        '''
        generate the 512-d embeddings of the texts
        '''

        with lock:
            inputs = self.processor(
                text=request.input, return_tensors="pt", padding=True).to(self.device)
            text_embeddings = self.model.get_text_features(**inputs)
            output = text_embeddings.cpu().detach().numpy().tobytes()

            return clip_pb2.EncodeResponse(embedding=output)

    def EncodeImage(self, request, context):
        '''
        generate the 512-d embeddings of the images
        '''
        with lock:
            image_data = request.image
            image = Image.open(io.BytesIO(image_data))

            inputs = self.processor(
                images=image, return_tensors="pt", padding=True).to(self.device)
            image_embeddings = self.model.get_image_features(**inputs)
            output = image_embeddings.cpu().detach().numpy().tobytes()

            return clip_pb2.EncodeResponse(embedding=output)

    def GetImageSpec(self, request, context):
        width = self.config.vision_config.image_size
        height = self.config.vision_config.image_size
        return clip_pb2.GetImageSpecResponse(width=width, height=height)


def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    clip_pb2_grpc.add_EncodingServiceServicer_to_server(
        EncodingServer(), server)
    server.add_insecure_port(os.getenv('CLIP_ADDRESS', '0.0.0.0:8173'))
    server.start()
    server.wait_for_termination()


if __name__ == '__main__':
    serve()
