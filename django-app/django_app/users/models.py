from django.db import models
from django.db.models.signals import post_save, post_delete
from django.dispatch import receiver

from grpc_commands.services.user_client import STUB
from pb import user_pb2


class ExternalUser(models.Model):
    name = models.CharField(max_length=100)
    email = models.CharField(max_length=100, unique=True)

    external_id = models.IntegerField(unique=True)
    from_grpc: bool = False

    class Meta:
        app_label = "users"


# method for updating
@receiver(post_save, sender=ExternalUser)
def update_user(sender, instance, **kwargs):
    if not instance.from_grpc and not kwargs.get('created'):
        user = user_pb2.User(
            id=instance.external_id, name=instance.name, email=instance.email
        )
        update_request = user_pb2.UpdateUserRequest(user=user)
        response = STUB.UpdateUser(update_request)
        print("Response", response)


@receiver(post_delete, sender=ExternalUser)
def delete_user(sender, instance, **kwargs):
    if not instance.from_grpc and not kwargs.get('created'):
        delete_request = user_pb2.DeleteUserRequest(id=instance.external_id)
        response = STUB.DeleteUser(delete_request)
        print("Response", response)
