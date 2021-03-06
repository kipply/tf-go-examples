package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"

	tfServing "tensorflow_serving/apis"

	tfTensorProto "github.com/tensorflow/tensorflow/tensorflow/go/core/framework/tensor_go_proto"
	tfTensorShapeProto "github.com/tensorflow/tensorflow/tensorflow/go/core/framework/tensor_shape_go_proto"
	tfTypesProto "github.com/tensorflow/tensorflow/tensorflow/go/core/framework/types_go_proto"
)

func main() {
	modelURL := "localhost:9000"
	modelMetadataRequest := &tfServing.GetModelMetadataRequest{
		ModelSpec: &tfServing.ModelSpec{
			Name: "model",
		},
		MetadataField: []string{"signature_def"},
	}

	conn, err := grpc.Dial(modelURL, grpc.WithInsecure())
	if err != nil {
		log.Fatalf(err.Error())
	}

	client := tfServing.NewPredictionServiceClient(conn)

	metadata, err := client.GetModelMetadata(context.Background(), modelMetadataRequest)
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Println(metadata)

	request := &tfServing.PredictRequest{
		ModelSpec: &tfServing.ModelSpec{
			Name:          "model",
			SignatureName: "sample",
		},
		Inputs: map[string]*tfTensorProto.TensorProto{
			"x": {
				Dtype: tfTypesProto.DataType_DT_INT64,
				TensorShape: &tfTensorShapeProto.TensorShapeProto{
					Dim: []*tfTensorShapeProto.TensorShapeProto_Dim{
						{
							Size: 1,
						},
						{
							Size: 5,
						},
					},
				},
				Int64Val: []int64{1, 2, 3, 4, 5},
			},
			"y": {
				Dtype: tfTypesProto.DataType_DT_INT32,
				TensorShape: &tfTensorShapeProto.TensorShapeProto{
					Dim: []*tfTensorShapeProto.TensorShapeProto_Dim{},
				},
				IntVal: []int32{1},
			},
			"z": {
				Dtype: tfTypesProto.DataType_DT_INT32,
				TensorShape: &tfTensorShapeProto.TensorShapeProto{
					Dim: []*tfTensorShapeProto.TensorShapeProto_Dim{},
				},
				IntVal: []int32{1},
			},
		},
	}
	inference, err := client.Predict(context.Background(), request)
	if err != nil {
		log.Fatalf(err.Error())
	}

	fmt.Println(inference)
}
