import grpc


class AuthInterceptor(grpc.ServerInterceptor):
    def __init__(self, valid_tokens):
        self.valid_tokens = valid_tokens

        def abort(ignored_request, context):
            context.abort(grpc.StatusCode.UNAUTHENTICATED, 'Invalid signature')

        self._abortion = grpc.unary_unary_rpc_method_handler(abort)

    def _validate_token(self, token):
        return token in self.valid_tokens

    def intercept_service(self, continuation, handler_call_details):
        metadata = dict(handler_call_details.invocation_metadata)
        token = metadata.get("authorization")

        if token is None and not self._validate_token(token):
            return self._abortion

        return continuation(handler_call_details)
