from django.db import models


class ExternalUser(models.Model):
    name = models.CharField(max_length=100)
    email = models.CharField(max_length=100, unique=True)
    external_id = models.IntegerField(unique=True)

    class Meta:
        app_label = 'users'
