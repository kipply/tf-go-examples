PATH_TO_TF=/Users/kipply/code/tensorflow/
PATH_TO_TF_SERVING=/Users/kipply/code/tfserving/

eval "protoc -I $PATH_TO_TF_SERVING -I $PATH_TO_TF --go_out=plugins=grpc:src/ $PATH_TO_TF_SERVING/tensorflow_serving/apis/*.proto"
eval "protoc -I $PATH_TO_TF_SERVING -I $PATH_TO_TF --go_out=plugins=grpc:src/ $PATH_TO_TF_SERVING/tensorflow_serving/config/*.proto"
eval "protoc -I $PATH_TO_TF_SERVING -I $PATH_TO_TF --go_out=plugins=grpc:src/ $PATH_TO_TF_SERVING/tensorflow_serving/util/*.proto"
eval "protoc -I $PATH_TO_TF_SERVING -I $PATH_TO_TF --go_out=plugins=grpc:src/ $PATH_TO_TF_SERVING/tensorflow_serving/core/*.proto"
eval "protoc -I $PATH_TO_TF_SERVING -I $PATH_TO_TF --go_out=plugins=grpc:src/ $PATH_TO_TF_SERVING/tensorflow_serving/sources/storage_path/*.proto"
eval "protoc -I $PATH_TO_TF_SERVING -I $PATH_TO_TF --go_out=plugins=grpc:src/ $PATH_TO_TF/tensorflow/core/framework/*.proto"
eval "protoc -I $PATH_TO_TF_SERVING -I $PATH_TO_TF --go_out=plugins=grpc:src/ $PATH_TO_TF/tensorflow/core/example/*.proto"
eval "protoc -I $PATH_TO_TF_SERVING -I $PATH_TO_TF --go_out=plugins=grpc:src/ $PATH_TO_TF/tensorflow/core/lib/core/*.proto"
eval "protoc -I $PATH_TO_TF_SERVING -I $PATH_TO_TF --go_out=plugins=grpc:src/ $PATH_TO_TF/tensorflow/core/protobuf/*.proto"
eval "protoc -I $PATH_TO_TF_SERVING -I $PATH_TO_TF --go_out=plugins=grpc:src/ $PATH_TO_TF/tensorflow/stream_executor/*.proto"

rm src/tensorflow_serving/apis/prediction_log.pb.go # causes an import loop
