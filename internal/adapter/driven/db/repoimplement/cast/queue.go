package cast

import (
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/moazedy/todo/internal/domain/model"
)

func ToMessageModel(in types.Message) (out model.Message) {
	if in.MessageId != nil {
		out.ID = *in.MessageId
	}

	if in.Body != nil {
		out.Body = *in.Body
	}

	return
}

func ToSliceOfModelMessages(in []types.Message) (out []model.Message) {
	for _, m := range in {
		out = append(out, ToMessageModel(m))
	}

	return
}
