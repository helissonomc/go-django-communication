import logging
from django.core.management.base import BaseCommand
from grpc_commands.services.user_service import server


class Command(BaseCommand):
    help = 'Launches Listener to get updates from Provider through Kafka'

    def handle(self, *args, **options):
        logging.basicConfig(level=logging.INFO)  # Set the logging level to INFO
        logging.getLogger().addHandler(logging.StreamHandler())
        server()
