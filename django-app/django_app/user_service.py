import grpc
import logging
from concurrent import futures
import user_pb2
import user_pb2_grpc


class UserService(user_pb2_grpc.UserServiceServicer):
    def CreateUser(self, request, context):
        # Implement logic to create user in Django
        # Return the created user
        print(f"Received request to create user: {request.user}")
        user = user_pb2.User(id=1, name=request.user.name, email=request.user.email)
        return user_pb2.CreateUserResponse(user=user)

    def UpdateUser(self, request, context):
        # Implement logic to update user in Django
        # Return the updated user
        user = user_pb2.User(id=request.user.id, name=request.user.name, email=request.user.email)
        return user_pb2.UpdateUserResponse(user=user)

    def DeleteUser(self, request, context):
        # Implement logic to delete user in Django
        # Return a response indicating success
        return user_pb2.DeleteUserResponse(success=True)


def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    user_pb2_grpc.add_UserServiceServicer_to_server(UserService(), server)
    server.add_insecure_port('[::]:50052')
    server.start()
    print("Server started, listening on 50052")
    server.wait_for_termination()


if __name__ == '__main__':
    serve()
    logging.basicConfig()
