import grpc
from concurrent import futures
import clip_pb2_grpc
import clip_pb2
import torch
from transformers import CLIPModel, CLIPProcessor

MODEL_ID = "openai/clip-vit-base-patch32"

class EncodingServer(clip_pb2_grpc.EncodingService):
    def __init__(self) -> None:
        self.device = "cuda" if torch.cuda.is_available() else "cpu"
        self.model = CLIPModel.from_pretrained(self.hf_model).to(self.device)
        self.processor = CLIPProcessor.from_pretrained(self.hf_model)
        self.logit_scale = self.model.logit_scale.item() if self.model.logit_scale.item() else 4.60517
        print("Model clip loaded", "device:", self.device)

    def EncodeText(self, request, context):
        '''
        generate the 512-d embeddings of the texts
        '''
        inputs = self.processor(text=request.text, return_tensors="pt", padding=True).to(self.device)
        text_embeddings = self.model.get_text_features(**inputs)
        output = text_embeddings.cpu().detach().numpy().tobytes()

        return clip_pb2.EncodeResponse(embedding=output)
    
    def EncodeImage(self, request, context):
        '''
        generate the 512-d embeddings of the images
        '''
        items = [request.image]
        inputs = self.processor(images=items, return_tensors="pt", padding=True).to(self.device)
        image_embeddings = self.model.get_image_features(**inputs)
        output = image_embeddings.cpu().detach().numpy()

        return clip_pb2.EncodeResponse(embedding=output)
    
    def GetImageSpec(self, request, context):
        return clip_pb2.GetImageSpecResponse(width=224, height=224)

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    clip_pb2_grpc.add_EncodingServiceServicer_to_server(EncodingServer(), server)
    server.add_insecure_port('[::]:8137')
    server.start()
    server.wait_for_termination()

if __name__ == '__main__':
    serve()