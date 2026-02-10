package codec

import (
	"github.com/boginskiy/GophKeeper/server/internal/model"
	"github.com/boginskiy/GophKeeper/server/internal/rpc"
	"github.com/boginskiy/GophKeeper/server/internal/utils"
)

type TextDecoderEncoder struct{}

func NewTextDecoderEncoder() *TextDecoderEncoder {
	return &TextDecoderEncoder{}
}

// >> Create

// DecoderCreateRequest
func (a *TextDecoderEncoder) DecoderCreateRequest(req *rpc.CreateRequest) (*model.Text, error) {
	return &model.Text{
		Name:         req.GetName(),
		Type:         req.GetType(),
		Content:      req.GetText(),
		Owner:        req.GetOwner(),
		ListActivate: req.GetListActivate(),
	}, nil
}

// EncodeCreateResponse
func (a *TextDecoderEncoder) EncodeCreateResponse(mod *model.Text) (*rpc.CreateResponse, error) {
	return &rpc.CreateResponse{
		Status:    "created",
		UpdatedAt: utils.ConvertDtStr(mod.UpdatedAt),
	}, nil
}

// >> Read

// DecoderReadRequest
func (a *TextDecoderEncoder) DecoderReadRequest(req *rpc.ReadRequest) (*model.Text, error) {
	return &model.Text{
		Name:  req.GetName(),
		Owner: req.GetOwner(),
	}, nil
}

// EncodeReadResponse
func (a *TextDecoderEncoder) EncodeReadResponse(mod *model.Text) (*rpc.ReadResponse, error) {
	return &rpc.ReadResponse{
		Status:    "read",
		Name:      mod.Name,
		Text:      mod.Content,
		UpdatedAt: utils.ConvertDtStr(mod.UpdatedAt),
	}, nil
}

// >> ReadAll

// DecoderReadAllRequest
func (a *TextDecoderEncoder) DecoderReadAllRequest(req *rpc.ReadAllRequest) (*model.Text, error) {
	return &model.Text{
		Type:  req.GetType(),
		Owner: req.GetOwner(),
	}, nil
}

// EncodeReadAllResponse
func (a *TextDecoderEncoder) EncodeReadAllResponse(mods []*model.Text) (*rpc.ReadAllResponse, error) {
	textResponses := make([]*rpc.TextResponse, len(mods))

	for i, text := range mods {
		textResponses[i] = &rpc.TextResponse{
			UpdatedAt: utils.ConvertDtStr(text.UpdatedAt),
			Text:      text.Content,
			Name:      text.Name,
		}
	}
	return &rpc.ReadAllResponse{
		Status:        "read",
		TextResponses: textResponses}, nil
}

// >> Delete

// DecoderDeleteRequest
func (a *TextDecoderEncoder) DecoderDeleteRequest(req *rpc.DeleteRequest) (*model.Text, error) {
	return &model.Text{
		Name:  req.GetName(),
		Owner: req.GetOwner(),
	}, nil
}

// EncodeDeleteResponse
func (a *TextDecoderEncoder) EncodeDeleteResponse(mod *model.Text) (*rpc.DeleteResponse, error) {
	return &rpc.DeleteResponse{
		Status: "deleted",
		Name:   mod.Name,
	}, nil
}
