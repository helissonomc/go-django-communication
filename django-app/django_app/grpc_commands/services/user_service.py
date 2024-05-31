import grpc
import logging
from concurrent import futures
from pb import user_pb2
from pb import user_pb2_grpc

from grpc_commands.services.auth_interceptor import AuthInterceptor
from users.models import ExternalUser


class UserService(user_pb2_grpc.UserServiceServicer):
    def CreateUser(self, request, context):
        # Implement logic to create user in Django
        # Return the created user
        user = user_pb2.User(
            id=request.user.id, name=request.user.name, email=request.user.email
        )
        external_user = ExternalUser
        external_user.from_grpc = True
        external_user.objects.create(
            external_id=user.id, name=user.name, email=user.email
        )

        response = user_pb2.CreateUserResponse(user=user)
        logging.info(response)
        return response

    def UpdateUser(self, request, context):
        # Implement logic to update user in Djangoxw
        # Return the updated user
        user = user_pb2.User(
            id=request.user.id, name=request.user.name, email=request.user.email
        )
        logging.info(user)
        ExternalUser.objects.filter(external_id=user.id).update(
            email=user.email,
            name=user.name,
        )
        return user_pb2.UpdateUserResponse(user=user)

    def DeleteUser(self, request, context):
        # Implement logic to delete user in Django
        # Return a response indicating success
        ExternalUser.objects.filter(external_id=request.id).delete()
        return user_pb2.DeleteUserResponse(success=True)


def server():
    valid_tokens = ["token_test", "your-valid-token2"]
    auth_interceptor = AuthInterceptor(valid_tokens)
    server = grpc.server(
        futures.ThreadPoolExecutor(max_workers=10), interceptors=(auth_interceptor,)
    )

    user_pb2_grpc.add_UserServiceServicer_to_server(UserService(), server)
    server.add_insecure_port("[::]:50052")
    server.start()
    print("Server started, listening on 50052")
    server.wait_for_termination()
