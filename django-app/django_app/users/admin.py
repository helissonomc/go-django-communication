from django.contrib import admin

from users.models import ExternalUser


@admin.register(ExternalUser)
class ExternalUserAdmin(admin.ModelAdmin):
    list_display = ("name", "email", "external_id")

    def has_add_permission(self, request):
        return False
