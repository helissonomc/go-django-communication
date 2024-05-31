from django.db import models
from django.db.models.signals import post_save
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
        create_request = user_pb2.UpdateUserRequest(user=user)
        response = STUB.UpdateUser(create_request)
        print("Response", response)
