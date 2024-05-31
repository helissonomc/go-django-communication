import grpc
from pb import user_pb2_grpc


channel = grpc.insecure_channel('go-grpc-server:50051')
STUB = user_pb2_grpc.UserServiceStub(channel)
